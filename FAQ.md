# FAQ


If you have any questions, please contact us:

* Email: [zack@hypercubeplatforms.com](mailto:zack@hypercubeplatforms.com)
* Gitter: [Join room](https://gitter.im/schollz/find)
* Github Issues: [Open an issue](https://github.com/schollz/find/issues/new)

###  What's the point of this?

**FIND** is a Framework for Internal Navigation and Discovery. Anything that you would do with a motion sensor you can do with FIND. Anything you can do with GPS information you can do with FIND information. Except here you get internal positioning so you could tell apart one table from another in a cafeteria, or one bookshelf from another in a library.

### How does it work?

The basis of this system is to catalog all the fingerprints about the
Wifi routers in the area (MAC addresses and signal values) and then
classify them according to their location. A typical set of fingerprints
froms some locations will look something like this (taken from the
library at Duke University):

![The distributions in Wifi RSSI signal as interpreted by an Android
device is shown for each router. Each router is differentiated by color.
Different locations give different distributions of router signals,
whether these locations differ by a several meters or the same location
between floors.](https://www.internalpositioning.com/guide/img/fingerprint_library.png)

### Can it replace motion sensors?

Yes...probably! Replacing motion sensors with FIND has the added benefit of automatically providing *user information* as well as the position/time information. The main disadvantage is that there is time delay of 1-5 seconds to register, so timing applications are not as good.


### Does it use a [Wi-Fi location database](https://en.wikipedia.org/wiki/Wi-Fi_positioning_system#Public_Wi-Fi_location_databases)?

No, there is no dependency on external resources like Wi-Fi location databases.

###  Doesn't this already exist?

We believe that our framework rivals and, for the most part, outperforms the currently existing frameworks for indoor geolocation and internal positioning. We certainly do not mean for you to take our word for this claim, though. To this end, we've decided to keep the project open-source so you can test it. We also invite you to make comparisons to current technologies:

If you are looking for a more **commercial, large-scale deployable application**, look at these up-and-coming solutions:

-   [MazeMap Indoor Navigation] - a Norway-based and Cisco-partnered enterprise that takes your CAD floor plans and generates a nice user-interface with similar indoor-positioning capabilities.
-   [Meridian Kits] - a SF and Portland based company (part of Aruba Networks) that offers specialized App SDK environments for building internal positioning systems into workplaces, businesses and hospitals
-   [MPact Platform] - Motorola is working on a internal positionings system that takes advantage of BlueTooth beacons and Wi-Fi for internal positioning for large applications

If you are looking for a **free single-user, single-home application**, perhaps you can find solutions with these apps:

-   [Google Maps Floor Plan Maker] - not sure how it works (and have not tested) but claims to be able to navigate within small businesses. Reviewed okay.
-   [WiFi Indoor Localization] - single-floor grid-based learning system that uses Wi-Fi to train on the names of mac addresses. In my hands it did not work well below 20ft resolution. Reviewed okay.
-   [Indoor Positioning] - Selective learning, not tested by me, but also grid-based. Not reviewed.
-   [BuildNGO - Indoor Navi] - Offers Android app that requires online service for uploading floor plans to their server and uses learning based on Max signal, may require Bluetooth as well.
-   [Wifarer] - Uses Beacons and WiFi for Indoor positioning, but trainable and limited to select museums. Reviewed well, but no training available.
-   [Indoor GPS] - Perfunctory application that trains on a route, instead of a location and offers SDK but still lots of work to be done. Reviewed okay.

  [MazeMap Indoor Navigation]: http://mazemap.com/what-it-is
  [Meridian Kits]: http://www.meridianapps.com
  [MPact Platform]: http://newsroom.motorolasolutions.com/Press-Releases/Communicate-to-Shoppers-at-the-Right-Time-with-First-of-its-Kind-Location-Based-Platform-from-Motor-49e1.aspx
  [Google Maps Floor Plan Maker]: https://play.google.com/store/apps/details?id=com.google.android.apps.insight.surveyor&hl=en
  [WiFi Indoor Localization]: https://play.google.com/store/apps/details?id=com.hfalan.wifilocalization&hl=en
  [Indoor Positioning]: https://play.google.com/store/apps/details?id=com.bombao.projetwifi&hl=en
  [BuildNGO - Indoor Navi]: https://play.google.com/store/apps/details?id=com.sails.buildngo&hl=en
  [Wifarer]: https://play.google.com/store/apps/details?id=com.wifarer.android&hl=en
  [Indoor GPS]: https://play.google.com/store/apps/details?id=com.ladiesman217.indoorgps&hl=en


###  Do I need to be on Wifi to submit fingerprints?

No, as long as you have a data connection. You also need to have Wifi enabled, otherwise you have no way of scanning wifi!

###  Can I run the server at home on the LAN connection?

Yes! You can setup your own server to host locally.

###  Can I track myself on a map?

Yes, soon!

###  Can I use an iPhone?

No. We currently do not support iPhone, and probably never will because they have locked developers from using the Wifi data.

###  Can I use an Android?

Yes! All Android devices are supported.

###  What is the minimum distance that can be resolved?

Less than 10 square feet. This depends highly on the surrounding location. It can be a few meters, possibly less for some locations.

###  Can you run the server on a Raspberry Pi?

Yes! Its been tested and runs great on a Raspberry Pi model B+.

### Can I help develop?

Yes! We host our code on [Github](https://github.com/schollz/find) and will accept Pull requests and Feature requests.
