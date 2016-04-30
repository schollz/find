# MQTT

# Setup (self-hosted servers)

If your not hosting, and just want to connect, read below.

Basically, right now you can only run `mosquitto` using a instance running from a configuration file created specifically by *FIND*. To get started, you'll first need the latest copy of `mosquitto`:

```bash
wget http://repo.mosquitto.org/debian/mosquitto-repo.gpg.key
sudo apt-key add mosquitto-repo.gpg.key
cd /etc/apt/sources.list.d/
sudo wget http://repo.mosquitto.org/debian/mosquitto-wheezy.list
sudo apt-get update
sudo apt-get install mosquitto-clients mosquitto
```

Then goto your FIND folder and create a file (in the future I'll have FIND do this automatically):

```
mkdir /path/to/find/mosquitto
touch /path/to/find/mosquitto/conf
```

Now, start `mosquitto` in the background:

```
mosquitto -c /path/to/find/mosquitto/conf -d
```

To use it with **FIND** you'll need the PID, so get that using

```bash
$ ps aux | grep mosquitto | grep conf
bitnami  PID  0.0  0.0  30968  1136 ?  Ss  09:19  0:00 mosquitto -c /path/to/find/mosquitto/conf -d
```

Now, you can startup **FIND**:

```bash
./find -mqtt ADDRESS:1883 -mqttadmin ADMIN -mqttadminpass ADMIN_PASS -mosquitto PID -p :PORT ADDRESS:PORT
```

The `ADDRESS` and `PORT` is the address and port your using for **FIND**. The `PID` is the `mosquitto` PID. The `ADMIN` and the `ADMIN_PASS` are your secret name and password to access read/write access to every MQTT channel. Make sure these are not simple enough to guess.

That's it!

# Client (MQTT connections)

## Register
To receive data from the **FIND** MQTT, follow these steps. First, register your group using the following:

```bash
curl -X PUT "https://ml.internalpositioning.com/mqtt?group=YOURGROUP"
```

where `YOURGROUP` is your group name. This command will tell **FIND** to add group level access to your own special MQTT channel. You'll receive a message like:

```javascript
{
    "message": "You have successfully set your password.",
    "password": "YOURPASSWORD",
    "success": true
}
```

The `password` is what you can use to access MQTT now. You can retrieve your password by using the same `curl` command. These passwords are completely random, and not hashed - so totally not guessable.

## Subscribing

First make sure to register. To subscribe to your channel to see current locations, simply use the topic `YOURGROUP/location/#`, e.g.:

```bash
mosquitto_sub -h ml.internalpositioning.com -u YOURGROUP -P YOURPASSWORD -t "YOURGROUP/location/#"
```

## Publishing Fingerprints

Currently, MQTT takes only a very specific type of fingerprint. Basically, to utilize the minimal MQTT byte size you have to compress the mac addresses and RSSI components.

To publish fingerprints, use the channel `YOURGROUP/track/USERNAME` for tracking or `YOURGROUP/learn/USERNAME/LOCATION` for learning. The body needs to be a multiple of 14 bytes where the first 12 bytes are the Mac address and the next 2 bytes is the RSSI value (absolute value). For example, if your detected routers are

```javascript
"ab:cd:ef:gf:ij:kl":-32
"mn:op:qr:st:uv:wx":-3
```

then you'll need to send the following as the body:

```
"abcdefgfijkl32mnopqrstuvwx 3"
```

Notice, that for absolute values < 10, you'll have to add the space. Here's an example Photon sketch for doing this:

```c
// This #include statement was automatically added by the Particle IDE.
#include "MQTT/MQTT.h"

// SWITCH
unsigned int SLEEP = 0;
unsigned int GREENLIGHT = 0;
unsigned int REDLIGHT = 0;
unsigned int nextTime = 0;
String group = "YOURGROUP";
String user = "YOURUSER";
String server = "ml.internalpositioning.com";
String password = "YOURPASSWORD"; // set with curl -X PUT "https://ml.internalpositioning.com/mqtt?group=YOURGROUP"
MQTT client("ml.internalpositioning.com", 1883, NULL );

void button_handler(system_event_t event, int duration, void* )
{
    if (!duration) { // just pressed
        RGB.control(true);
        if (SLEEP == 0) {
            RGB.color(0,255,0);
            GREENLIGHT = 1;
            SLEEP = 1; // sleep mode on
        } else {
            RGB.color(255,0,0);
            REDLIGHT = 1;
            SLEEP = 0; // sleep mode off
        }
    }
    else {    // just released
        RGB.control(false);
        GREENLIGHT = 0;
        REDLIGHT = 0;
    }
}

void setup() {
    RGB.control(false);

    // connect to the server
    client.connect(server,group,password);
    System.on(button_status, button_handler);
}

void loop() {

    // SWITCH
    if (REDLIGHT == 1) {
        RGB.color(255,0,0);
    }
    if (GREENLIGHT == 1) {
        RGB.color(0,255,0);
    }

    if (nextTime > millis()) {
        return;
    }


    if (client.isConnected()) {
        String body;
        body = "";
        WiFiAccessPoint aps[20];
        int found = WiFi.scan(aps, 20);
        for (int i=0; i<found; i++) {
            WiFiAccessPoint& ap = aps[i];
            char mac[17];
            sprintf(mac,"%02x%02x%02x%02x%02x%02x",
             ap.bssid[0] & 0xff, ap.bssid[1] & 0xff, ap.bssid[2] & 0xff,
             ap.bssid[3] & 0xff, ap.bssid[4] & 0xff, ap.bssid[5] & 0xff);
            body = body + mac;
            float f = -1*ap.rssi;
            if (f > 100) {
                f = 99;
            }
            String sf(f, 0);
            if (f < 10) {
                body = body + " " + sf;
            } else {
                body = body + sf;
            }
        }

        client.publish(group + "/track/" + user,body); // Use this for tracking
        // client.publish(group + "/learn/" + user + "/desk",body); // Use this for learning

        nextTime = millis() + 2000; // sends response every 5 seconds  (2 sec delay + ~3 sec for gathering signals)

        // // SWITCH
        if (SLEEP == 1) {
            System.sleep(5);
        }
        client.loop();

    } else {
        client.connect(server,group,password);
    }
}

```
