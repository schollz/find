# Integration with Particle Photon

There are two ways you can integrate a Particle Photon into find. First, you can simply place the Photon in listening mode, and if you have several you can put them in different locations. These will especially help with geolocation.

Second, you can use the Particle Photon as a tracking device using the following code. The code allows for a "sleep" mode that can be activated by pressing "Setup" once, and turned off by pressing "Setup" twice. This mode allows more battery life.
Here is my `.ino` file:

```c
#include "application.h"
#include "HttpClient/HttpClient.h"

/**
* Declaring the variables.
*/
unsigned int nextTime = 0;    // Next time to contact the server
HttpClient http;

// Headers currently need to be set at init, useful for API keys etc.
http_header_t headers[] = {
    { "Content-Type", "application/json" },
    { "Accept" , "application/json" },
    { NULL, NULL } // NOTE: Always terminate headers will NULL
};

http_request_t request;
http_response_t response;

// // SWITCH
// unsigned int SLEEP = 0;
// void button_handler(system_event_t event, int duration, void* )
// {
//     if (!duration) { // just pressed
//         RGB.control(true);
//         if (SLEEP == 0) {
//             RGB.color(255,0,0);
//             SLEEP = 1; // sleep mode on
//         } else if (SLEEP == 1) {
//             RGB.color(0,255,0);
//             SLEEP = 2; // undefined mode
//         } else {
//             RGB.color(0,0,255);
//             SLEEP = 0;
//         }
//     }
//     else {    // just released
//         RGB.control(false);
//     }
// }

void setup() {
    Serial.begin(9600);
    // // SWITCH
    // System.on(button_status, button_handler);
}

void loop() {
    
    // // SWITCH
    // if (SLEEP == 1) {
    //     RGB.color(255,0,0);
    // }
    // if (SLEEP == 2) {
    //     RGB.color(0,255,0);
    // }

    if (nextTime > millis()) {
        return;
    }

    // // DEBUGGING
    // Serial.println();
    // Serial.println("Application>\tStart of Loop.");

    request.hostname = "ml2.internalpositioning.com";
    request.port = 80;
    request.path = "/track";


    request.body = "{\"group\":\"YOURGROUP\",\"username\":\"YOURUSERNAME\",\"location\":\"YOURLOCATION\",\"wifi-fingerprint\":[";
    WiFiAccessPoint aps[20];
    int found = WiFi.scan(aps, 20);
    for (int i=0; i<found; i++) {
        WiFiAccessPoint& ap = aps[i];
        char mac[17];
        sprintf(mac,"%02x:%02x:%02x:%02x:%02x:%02x",
         ap.bssid[0] & 0xff, ap.bssid[1] & 0xff, ap.bssid[2] & 0xff,
         ap.bssid[3] & 0xff, ap.bssid[4] & 0xff, ap.bssid[5] & 0xff);
        request.body = request.body + "{\"mac\":\"" + mac + "\",";
        float f = ap.rssi;
        String sf(f, 0);
        request.body = request.body + "\"rssi\":" + sf + "}";
        
        if (i < found -1 ) {
            request.body = request.body + ",";
        }
    }
    request.body = request.body + "]}";

    http.post(request, response, headers);
    
    // // DEBUGGING
    // Serial.println("Fingerprint:");
    // Serial.println(request.body);
    // Serial.print("Application>\tResponse status: ");
    // Serial.println(response.status);
    // Serial.print("Application>\tHTTP Response Body: ");
    // Serial.println(response.body);
    
    nextTime = millis() + 2000; // sends response every 5 seconds  (2 sec delay + ~3 sec for gathering signals)
    
    // // SWITCH
    // if (SLEEP == 1) {
    //     System.sleep(3);
    // } else {
    //     delay(3000);
    // }
}
```

### Things to keep in mind

- You must use HTTP, not HTTPS. That's why the server is set to `ml2.internalpositioning.com`
- You can not flash from WiFi is the board is in sleep mode. Thats what the button is for. If this fails, you can reset by unplugging, holding down "Setup" and then plugging in while holding down "Setup." Then link up the Photon like you did from the beginning.
- The Photon ESP chip sees fewer macs than a Android does, probably because of the antenna. Thus, its best to not use platform-specific information and you should set the mixins to `0` on the server by using `curl https://ml.internalpositioning.com/mixin?group=X&mixin=0`.


### Benchmarking

Using a `nextTime` of `+2000ms` without System sleep gave a runtime of 13 hours, 32 minutes - ~813 minutes. Using fully charged [battery](http://www.insigniaproducts.com/products/computer-speakers-accessories/NS-MB2601.html) with 2600 mAh. This means, the sketch takes about 192 mA.
