# FAQ


If you have any questions, please contact:

* Email: [zack@hypercubeplatforms.com](mailto:zack@hypercubeplatforms.com)
* Twitter: [@zack_118](https://twitter.com/intent/tweet?screen_name=zack_118)
* Gitter: [Join room](https://gitter.im/schollz/find)
* Github Issues: [Open an issue](https://github.com/schollz/find/issues/new)

<br><br><br><br>


###  Can I use an iPhone?

**No.** We currently do not support iPhone. Unfortunately, the information about the WiFi scanning has to come from the use of the [`Apple80211` library](https://stackoverflow.com/questions/9684341/iphone-get-a-list-of-all-ssids-without-private-library/9684945#9684945). This is private library which means that [a user would have to jail break their device in order to use it](https://stackoverflow.com/questions/6341547/ios-can-i-manually-associate-wifi-network-with-geographic-location/6341893#6341893). We do not want to distribute an app that would require users to jailbreak their phones, so we will have to avoid developing for iOS until Apple removes this restriction. Sorry!

<br><br><br><br>


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


  <br><br><br><br>


###  What's the point of this?

**The point is to eventually incorporate FIND into home automation.** **FIND** can replace motion sensors to provide positional and user-specific information. Anything that you would do with a motion sensor you can do with **FIND**. Anything you can do with GPS information you can do with **FIND** information. Except here you get internal positioning so you could tell apart one table from another in a cafeteria, or one bookshelf from another in a library.

 As Wi-Fi singleboard computers get smaller and smartphones become more ubiquitous there will be more and more opportunities to harness WiFi signals into something useful for other applications.

 <br><br><br><br>


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

<br><br><br><br>


### But, really, how does it work?

**FIND** classifies locations using two methods.

#### Method 1

**FIND** uses a very simple, but very effective machine learning formulation known as [a Naive-Bayes classifier](https://en.wikipedia.org/wiki/Naive_Bayes_classifier) for **Method 1**.
The formula used to calculate probabilities of being in a location \( loc_x \) given your fingerprinting device sees \( max_y \) with signal strength \( z_y \) is

$$P ( loc_x | mac_y = z_y ) = \prod_{y=1}^{n} \left[  \frac{ P ( mac_y = z_y | loc_x  ) P(loc_x) }{ P (  mac_y = z_y ) } \right].$$

which is the product over all macs seen in a single fingerprint. The goal is to find the location \( x \) that maximizes the probability \( P ( loc_x | mac_y = z) \) for all macs \( 1..n \).

The probability \( P ( mac_y = z_y | loc_x  ) \) comes from calculating the distributions from the learned locations of the typical signals that each mac address gives out in a given location. Similarly, the term \( P ( mac_y = z_y ) \) comes from calculating the distribution of the mac having that signal, regardless of the location. The term \( P(loc_x) \) is simply the probability of being in that location which is assumed to be uniform so it is equal to \( 1/M \) where \( M \) is the total number of locations.

This *works*, but not as well as I thought it would. The reason is that there is information not being utilized about where macs do not have any signal, which can be addressed in the following implementation, known as **Naive-Bayes - Extended form (assuming exhaustive possibilities)**:

$$P ( loc_x | mac_y = z_y ) = \prod_{y=1}^{n} \left[  \frac{ P ( mac_y = z_y | loc_x  ) P(loc_x) }{ P ( mac_y = z_y | loc_x  ) P(loc_x) + P ( mac_y = z_y | \neg loc_x  ) P( \neg loc_x)} \right]$$

which is identical to above except now we calculate the probability distribution of the signal of a each mac in the location that is not x,  \( P ( mac_y = z_y | \neg loc_x  ) \). Also there is also a prior on the probability of being in a location that is not \( loc_x \), \( P( \neg loc_x) = 1 - 1/M \) where \( M \) is the total number of possible locations.

This form assumes that the locations are *mutually
exclusive* and *exhaustive* - which they are. You cannot be in more than one location at once and the only possible locations are those that are known.

**Practical considerations:**

*   For practical computational reasons, the logarithm of the \(P ( loc_x | mac_y = z_y )\) is calculated so that product becomes a sum.
*   Issue: If the numerator on the right side of the equation is ever zero, it
    cancels out that location completely (0 x anything = 0). This could
    happen if there is ever a mac address in the fingerprint that hasn't
    been learned. This could happen if a new router appears, or if
    someone picks up another stray router somehow. This shouldn't happen
    normally, but to prevent this case, perhaps all \(P ( mac_y = z_y | \neg loc_x  )\)
    distributions should be initialized with some small *non-zero* value.
*   The user sends a fingerprint of a small subset of routers, which have
    signal values associated with them. Since this equation iterates over
    ALL routers in ALL locations, for the routers that were not specifies
    (i.e. the routers of other locations) those should have signal -100.
    During the prior calculation, then, it is important to initialize all
    the distributions of rooms with nonexistent routers with a uniform
    function at signal ``absentee``.

#### Method 2

The second method is uses the *prevalence* of a specific Mac address in a single fingerprint. A single fingerprint is composed of all the Mac addresses and signal values that were detected in that read. In certain locations, you can more frequently sample certain routers than other routers. Thus, a very simple method of classification can use the frequency of "seeing" a Mac address as a basis.

**Practical considerations:**

* This calculation is *Antenna based* since different devices will be able to "see" more or less routers based on their antenna type. Thus, this type of calculation will often be used for single-device learning.


#### Mixin

There is a `mixin` variable that calculates the linear combination of Method 1 and 2 that results in the optimal parameter set.

**Practical considerations:**

* This `mixin` parameter [can be set manually](/#put-mixin), in the case that you are sure you don't want to use antenna specific information (Method 2).

<br><br><br><br>


### Does the smartphone version use up the battery quickly?

**No.** One important benefit of using WiFi-based technologies as they are relatively cheap sensors in the smartphone/computer. They are cheap in a monetary sense, as well as a power supply sense. Normally, a GPS sensor that is probed every 2 seconds will quickly drain your battery. Probing the WiFi every 2 seconds will take hours to drain your battery and is less taxing than many apps.

<br><br><br><br>

### Can it replace motion sensors?

**Yes...probably.** Replacing motion sensors with FIND has the added benefit of automatically providing *user information* as well as the position/time information. The main disadvantage is that there is time delay of 1-5 seconds to register, so timing applications are not as good. It is possible to increase the scan speed to accomplish better timing applications but it will drain the life of the battery faster.

<br><br><br><br>

### Does it use a [Wi-Fi location database](https://en.wikipedia.org/wiki/Wi-Fi_positioning_system#Public_Wi-Fi_location_databases)?

**No.** There is no dependency on external resources like Wi-Fi location databases. However, these type of databases can add additional information that might be worthwhile to explore to also integrate into **FIND**.

<br><br><br><br>

###  Do I need to be on Wifi to submit fingerprints?

**Yes, unless you have a data connection**. You also need to have Wifi enabled, otherwise you have no way of scanning wifi!

<br><br><br><br>


###  Can I run the server at home on the LAN connection?

**Yes.** You can setup your own server to host locally. Setting up your server can be done by [building the program yourself](https://github.com/schollz/find#setup-optional) or simply [downloading the latest prebuilt version](https://github.com/schollz/find/releases/tag/v2.0) for you OS.

<br><br><br><br>


###  Can I track myself on a map?

**Not yet**. This is something we would like to implement and we are working on. There is more information about our longterm roadmap [here](https://www.internalpositioning.com/).


<br><br><br><br>


###  Can I use an Android?

**Yes.** All Android devices are supported. You can [download the app from Google Play](https://play.google.com/store/apps/details?id=com.hcp.find) or [use the source code to build the app yourself](https://github.com/schollz/find/tree/android).

<br><br><br><br>


###  What is the minimum distance that can be resolved?

**It depends.** This system harnesses the available WiFi routers. If you have very few WiFi routers in the vicinity (i.e. <3 in 50 meters) then your resolution will suffer. Otherwise, you can typically get less than 10 square feet in location resolution.

<br><br><br><br>


###  Can you run the server on a Raspberry Pi?

**Yes.** Its been tested and runs great on a Raspberry Pi model B+, and model 3. Also, all the releases include [a ARM release for running on a Raspberry pi](https://github.com/schollz/find/releases).

<br><br><br><br>


### Is there a good minimum time to leave the app to train a location?

<br><br><br><br>


### Can it pick up locations between floors?

**Yes.** Yes it will pick up floors no problem. Floors tend to attenuate the signal, so there is a noticeable difference when you are in the same position, but on different floors. For example, check this out [this graphic](https://camo.githubusercontent.com/aa8ce49f332c0d1dcf0baa58c9a1d59672f8ce22/68747470733a2f2f7777772e696e7465726e616c706f736974696f6e696e672e636f6d2f67756964652f696d672f66696e6765727072696e745f6c6962726172792e706e67). Location x54 and location x49 are exactly the same place, but different floors, in a library. The blue signals are much more attenuated in x49 than x54 and thus are distinguished.

<br><br><br><br>


### What is a good amount of time to train a location?

**2 to 5 minutes**. Optimally you want to send ~100 pieces of information to the server. It transmits about 20 per minute, so you should give it some time to train well.

<br><br><br><br>


### Can I help develop?

**Yes!** We host our code on [Github](https://github.com/schollz/find) and will accept Pull requests and Feature requests.

<br><br><br><br>


### Does it work with [home-assistant.io](https://home-assistant.io/)?

**Yes.** See [here](https://community.home-assistant.io/t/anyone-seen-this-find-internal-positioning/772/2?u=schollz) for the discussion on how to use it with home-assistant.io.
