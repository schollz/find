# Integration with Particle Photon

The Particle Photon is a $20 device that can be used as a tracking device using the following code. The code allows for a "sleep" mode that can be activated by pressing "Setup" once, and turned off by pressing "Setup" twice. This mode allows more battery life.

You can access the **FIND** program using HTTP or MQTT. I've included example sketches that can be used for each.


### HTTP-based messaging

[Source code](https://gist.github.com/schollz/6077b4c64cf488e89856ed76c0f8a7d2).

Some notes:

- You must use HTTP, not HTTPS. That's why the server is set to `ml2.internalpositioning.com`
- You can not flash from WiFi is the board is in sleep mode. Thats what the button is for. If this fails, you can reset by unplugging, holding down "Setup" and then plugging in while holding down "Setup." Then link up the Photon like you did from the beginning.
- The Photon ESP chip sees fewer macs than a Android does, probably because of the antenna. Thus, its best to not use platform-specific information and you should set the mixins to `0` on the server by using `curl https://ml.internalpositioning.com/mixin?group=X&mixin=0`.


### MQTT-based messaging

[Source code](https://gist.github.com/schollz/d388b0c0ed1bb3b604eba6b7154a49d1).

This version uses significantly less bandwidth than the HTTP version. It takes a little more configuration, [see the MQTT documentation](/mqtt/) for how to get started with MQTT.

### Benchmarking

Using a `nextTime` of `+2000ms` it uses 147 mA.

Using a `nextTime` of `+5000ms`, it uses 113 mA.

Using a `nextTime` of `+10000ms`, with SLEEP ACTIVATED, it uses 81 mA

Using a `nextTime` of `+60000ms`, with SLEEP ACTIVATED, it uses 57 mA


