# MQTT

# Setup (Self-hosted servers)

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
