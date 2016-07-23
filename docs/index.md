[![](https://raw.githubusercontent.com/schollz/find/master/static/splash.gif)](https://www.internalpositioning.com/)



**Keywords**: indoor GPS, WiFi positioning, indoor mapping, indoor navigation, indoor positioning

# About

**The Framework for Internal Navigation and Discovery** (_FIND_) allows you to use your (Android) smartphone or WiFi-enabled computer (laptop or Raspberry Pi or etc.) to determine your position within your home or office. You can easily use this system in place of motion sensors as its resolution will allow your phone to distinguish whether you are in the living room, the kitchen or the bedroom, etc. The position information can then be used in a variety of ways including home automation, way-finding, or tracking!
<blockquote>Simply put, FIND will allow you to replace <em>tons</em> of motion sensors with a <em>single</em> smartphone!</blockquote>

The system is built on two main components - [a server](https://ml.internalpositioning.com/) and a fingerprinting device. The fingerprinting device ([computer program](https://github.com/schollz/find/releases/tag/v0.4client) or [android app](https://play.google.com/store/apps/details?id=com.hcp.find)) sends the specified data to the machine learning server which stores the fingerprints and analyzes them. It then returns the result to the device and stores the result on the server for accessing via a web browser or triggering via hooks.


# Quickstart

If you'd like to install things yourself, see [Setup documentation](https://doc.internalpositioning.com/setup/). You don't need to do this to try FIND, though.

## 1. Download the software

**Android users:** [download the current version of the app](https://play.google.com/store/apps/details?id=com.hcp.find). *Sorry iPhone users but  [the Apple store prevents apps that access WiFi information](https://github.com/schollz/find/blob/master/FAQ.md#can-i-use-an-iphone), so I will be unable to release a iPhone version.*

**Computer users:** you can [download the current version of the fingerprinting program](https://github.com/schollz/find/releases/tag/v0.4client), available for Rasbperry Pi, OSX, Linux, and Windows.


## 2. Gather fingerprint data

First, to get started using **FIND** you will need to gather fingerprint data in your locations.

**Android users:** When you start up the app you will be asked for a username (enter whatever you want) and you'll be assigned a unique group name. Simply click "Learn" and you'll be prompted for a location name. After you enter a location, the app will connect to the server and then submit fingerprints.

<center>
<img src="https://i.imgur.com/fbcYom5.png" width="200px" />
<img src="https://i.imgur.com/Ab9eXIk.png" width="200px" />
</center>
<br>


**Computer users:** To start learning locations simply use `./fingerprint -e`.



## 3. Track yourself

Once you've collected data in a few locations, you can track yourself.

**Android users:** Just press the "Track" button when you're ready to track.

**Computer users:** Type in `./fingerprint` to start tracking yourself.
