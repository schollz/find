# <img src="https://www.internalpositioning.com/guide/img/wifi-marker-darkgrey-small.png" width="30px" /> FIND




[![Join the chat at https://gitter.im/schollz/find](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/schollz/find?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) [![Version 2.1](https://img.shields.io/badge/version-2.1-brightgreen.svg)](https://www.internalpositioning.com/guide/development/)
[![FIND documentation](https://img.shields.io/badge/find-documentation-blue.svg)](https://doc.internalpositioning.com/)
[![Go Report Card](https://goreportcard.com/badge/github.com/schollz/find)](https://goreportcard.com/report/github.com/schollz/find) ![Coverage](https://img.shields.io/badge/coverage-34%25-orange.svg)
 [![Donate](https://img.shields.io/badge/donate-$1-brown.svg)](https://www.paypal.me/ZackScholl/1.00)

<a href="https://www.internalpositioning.com/"><img src="https://raw.githubusercontent.com/schollz/find/master/static/splash.gif"></a>

**Keywords**: indoor GPS, WiFi positioning, indoor mapping, indoor navigation, indoor positioning

# About

**The Framework for Internal Navigation and Discovery** (_FIND_) allows you to use your (Android) smartphone or WiFi-enabled computer (laptop or Raspberry Pi or etc.) to determine your position within your home or office. You can easily use this system in place of motion sensors as its resolution will allow your phone to distinguish whether you are in the living room, the kitchen or the bedroom, etc. The position information can then be used in a variety of ways including home automation, way-finding, or tracking!
<blockquote>Simply put, FIND will allow you to replace <em>tons</em> of motion sensors with a <em>single</em> smartphone!</blockquote>

The system is built on two main components - [a server](https://ml.internalpositioning.com/) and a fingerprinting device. The fingerprinting device ([computer program](https://github.com/schollz/find/releases/tag/v0.4client) or [android app](https://play.google.com/store/apps/details?id=com.hcp.find)) sends the specified data to the machine learning server which stores the fingerprints and analyzes them. It then returns the result to the device and stores the result on the server for accessing via a web browser or triggering via hooks.

**FAQ (abbreviated)**:
- How does it work? [It uses already available WiFi information to classify locations](https://github.com/schollz/find/blob/master/FAQ.md#how-does-it-work).
- Doesn't this already exist? [Yes, kinda](https://github.com/schollz/find/blob/master/FAQ.md#doesnt-this-already-exist).
- What's the point? This short piece of code can do [home automation](https://github.com/schollz/find/blob/master/FAQ.md#can-i-run-the-server-at-home-on-the-lan-connection) and [replace motion sensors](https://github.com/schollz/find/blob/master/FAQ.md#can-it-replace-motion-sensors) and  [more](https://github.com/schollz/find/blob/master/FAQ.md#whats-the-point-of-this).
- Can I use an iPhone? [Sorry, no](https://github.com/schollz/find/blob/master/FAQ.md#can-i-use-an-iphone).
- Does it work on a Raspberry Pi? [Yes](https://github.com/schollz/find/blob/master/FAQ.md#can-you-run-the-server-on-a-raspberry-pi).
- Does it work with [home-assistant.io](https://home-assistant.io/)? [Yes](https://github.com/schollz/find/blob/master/FAQ.md#does-it-work-with-home-assistantio).
- Can I help? [Yes, please](https://github.com/schollz/find/issues).
- How do I try it? It's easy. Just [download and run](https://github.com/schollz/find#usage).

More questions? See the [unabbreviated FAQ](https://github.com/schollz/find/blob/master/FAQ.md).


# Quickstart

If you'd like to install things yourself, see [INSTALL.md](https://github.com/schollz/find/blob/master/INSTALL.md). You don't need to do this to try FIND, though.

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


# More information

See the documentation at https://doc.internalpositioning.com.

### Donate

Like this? Help me keep it alive [by donating $5](https://www.paypal.me/ZackScholl/5.00) to [pay for server costs](http://rpiai.com/donate/).
