# FAQ


If you have any questions, please contact:

* Email: [zack@hypercubeplatforms.com](mailto:zack@hypercubeplatforms.com)
* Twitter: [@zack_118](https://twitter.com/intent/tweet?screen_name=zack_118)
* Gitter: [Join room](https://gitter.im/schollz/find)
* Github Issues: [Open an issue](https://github.com/schollz/find/issues/new)

###  Can I use an iPhone?

**No.** We currently do not support iPhone. Unfortunately, the information about the WiFi scanning has to come from the use of the [`Apple80211` library](https://stackoverflow.com/questions/9684341/iphone-get-a-list-of-all-ssids-without-private-library/9684945#9684945). This is private library which means that [a user would have to jail break their device in order to use it](https://stackoverflow.com/questions/6341547/ios-can-i-manually-associate-wifi-network-with-geographic-location/6341893#6341893). We do not want to distribute an app that would require users to jailbreak their phones, so we will have to avoid developing for iOS until Apple removes this restriction. Sorry!

###  Doesn't this already exist?

**Yes - but not satisfyingly.** Most solutions are not open-source, or they require external hardware (beacons, etc.), or they are expensive, or they just don't work very well. But don't take my word for it, try it yourself. Here are some of the programs I found that are similar:

If you are looking for a more **commercial, large-scale deployable application**, look at these up-and-coming solutions:

-   [MazeMap Indoor Navigation] - a Norway-based and Cisco-partnered enterprise that takes your CAD floor plans and generates a nice user-interface with similar indoor-positioning capabilities.
-   [Meridian Kits] - a SF and Portland based company (part of Aruba Networks) that offers specialized App SDK environments for building internal positioning systems into workplaces, businesses and hospitals
-   [MPact Platform] - Motorola is working on a internal positioning system that takes advantage of BlueTooth beacons and Wi-Fi for internal positioning for large applications

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


###  What's the point of this?

**The point is to eventually incorporate FIND into home automation.** **FIND** can replace motion sensors to provide positional and user-specific information. Anything that you would do with a motion sensor you can do with **FIND**. Anything you can do with GPS information you can do with **FIND** information. Except here you get internal positioning so you could tell apart one table from another in a cafeteria, or one bookshelf from another in a library.

 As Wi-Fi singleboard computers get smaller and smartphones become more ubiquitous there will be more and more opportunities to harness WiFi signals into something useful for other applications.

### How does it work?

**The basis of this system is to catalog all the fingerprints about the
Wifi routers in the area (MAC addresses and signal values) and then
classify them according to their location.** Take a look at a typical set of fingerprints
 (taken from the library at Duke University):

![Distributions](https://www.internalpositioning.com/guide/img/fingerprint_library.png)

The distributions in Wifi RSSI signal as interpreted by an Android
device is shown for each router. Each router is differentiated by color.
Different locations give different distributions of router signals,
whether these locations differ by a several meters or the same location
between floors.

**FIND** works by taking these differences between the WiFi data from different places to *classify* a location. Positioning is accomplished by first learning the distributions of WiFi signals for a given location and then classifying it during tracking. Learning only takes ~10 minutes and will last almost indefinitely. The WiFi fingerprints are also the same across all devices so that learning using one device is guaranteed to work across all devices.

### Does the smartphone version use up the battery quickly?

**No.** One important benefit of using WiFi-based technologies as they are relatively cheap sensors in the smartphone/computer. They are cheap in a monetary sense, as well as a power supply sense. Normally, a GPS sensor that is probed every 2 seconds will quickly drain your battery. Probing the WiFi every 2 seconds will take hours to drain your battery and is less taxing than many apps.

### Can it replace motion sensors?

**Yes...probably.** Replacing motion sensors with FIND has the added benefit of automatically providing *user information* as well as the position/time information. The main disadvantage is that there is time delay of 1-5 seconds to register, so timing applications are not as good. It is possible to increase the scan speed to accomplish better timing applications but it will drain the life of the battery faster.


### Does it use a [Wi-Fi location database](https://en.wikipedia.org/wiki/Wi-Fi_positioning_system#Public_Wi-Fi_location_databases)?

**No.** There is no dependency on external resources like Wi-Fi location databases. However, these type of databases can add additional information that might be worthwhile to explore to also integrate into **FIND**.


###  Do I need to be on Wifi to submit fingerprints?

**Yes, unless you have a data connection**. You also need to have Wifi enabled, otherwise you have no way of scanning wifi!

###  Can I run the server at home on the LAN connection?

**Yes.** You can setup your own server to host locally. Setting up your server can be done by [building the program yourself](https://github.com/schollz/find#setup-optional) or simply [downloading the latest prebuilt version](https://github.com/schollz/find/releases/tag/v2.0) for you OS.

###  Can I track myself on a map?

**Not yet**. This is something we would like to implement and we are working on. There is more information about our longterm roadmap [here](https://www.internalpositioning.com/).


###  Can I use an Android?

**Yes.** All Android devices are supported. You can [download the app from Google Play](https://play.google.com/store/apps/details?id=com.hcp.find) or [use the source code to build the app yourself](https://github.com/schollz/find/tree/android).

###  What is the minimum distance that can be resolved?

**It depends.** This system harnesses the available WiFi routers. If you have very few WiFi routers in the vicinity (i.e. <3 in 50 meters) then your resolution will suffer. Otherwise, you can typically get less than 10 square feet in location resolution.

###  Can you run the server on a Raspberry Pi?

**Yes.** Its been tested and runs great on a Raspberry Pi model B+, and model 3. Also, all the releases include [a ARM release for running on a Raspberry pi](https://github.com/schollz/find/releases).

### Is there a good minimum time to leave the app to train a location?

### Can it pick up locations between floors?

**Yes.** Yes it will pick up floors no problem. Floors tend to attenuate the signal, so there is a noticeable difference when you are in the same position, but on different floors. For example, check this out [this graphic](https://camo.githubusercontent.com/aa8ce49f332c0d1dcf0baa58c9a1d59672f8ce22/68747470733a2f2f7777772e696e7465726e616c706f736974696f6e696e672e636f6d2f67756964652f696d672f66696e6765727072696e745f6c6962726172792e706e67). Location x54 and location x49 are exactly the same place, but different floors, in a library. The blue signals are much more attenuated in x49 than x54 and thus are distinguished.

### What is a good amount of time to train a location?

**2 to 5 minutes**. Optimally you want to send ~100 pieces of information to the server. It transmits about 20 per minute, so you should give it some time to train well.

### Can I help develop?

**Yes!** We host our code on [Github](https://github.com/schollz/find) and will accept Pull requests and Feature requests.

### Does it work with [home-assistant.io](https://home-assistant.io/)?

**Yes.** See [here](https://community.home-assistant.io/t/anyone-seen-this-find-internal-positioning/772/2?u=schollz) for the discussion on how to use it with home-assistant.io.
