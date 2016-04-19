# <img src="https://www.internalpositioning.com/guide/img/wifi-marker-darkgrey-small.png" width="30px" /> FIND




[![Join the chat at https://gitter.im/schollz/find](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/schollz/find?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) [![Version 2.0](https://img.shields.io/badge/version-2.0-brightgreen.svg)](https://www.internalpositioning.com/guide/development/)
[![Go Report Card](https://goreportcard.com/badge/github.com/schollz/find)](https://goreportcard.com/report/github.com/schollz/find) ![Coverage](https://img.shields.io/badge/coverage-21%25-orange.svg) [![Donate](https://img.shields.io/badge/donate-$1-brown.svg)](https://www.paypal.me/ZackScholl/1.00)

<a href="https://www.internalpositioning.com/"><img src="https://raw.githubusercontent.com/schollz/find/master/static/splash.gif"></a>

*FIND is now 100% open-source.*

*Note for past users, the codebase has been completely rewritten in Golang so [things are ~100x faster](https://github.com/schollz/find/blob/master/BENCHMARKS.md#python-vs-go), smarter, and more secure. The [Python3 version will stay available](https://github.com/schollz/find/tree/python3), but it is no longer supported. Development will now be focused on this version.*

**Keywords**: indoor GPS, WiFi positioning, indoor mapping, indoor navigation, indoor positioning

# About

**The Framework for Internal Navigation and Discovery** (_FIND_) allows you to use your (Android) smartphone or WiFi-enabled computer (laptop or Raspberry Pi or etc.) to determine your position within your home or office. You can easily use this system in place of motion sensors as its resolution will allow your phone to distinguish whether you are in the living room, the kitchen or the bedroom, etc. The position information can then be used in a variety of ways including home automation, way-finding, or tracking!
<blockquote>Simply put, FIND will allow you to replace <em>tons</em> of motion sensors with a <em>single</em> smartphone!</blockquote>

The system is built on two main components - [a server](https://ml.internalpositioning.com/) and a fingerprinting device. The fingerprinting device ([computer program](https://github.com/schollz/find/releases/tag/v0.3client) or [android app](https://play.google.com/store/apps/details?id=com.hcp.find)) sends the specified data to the machine learning server which stores the fingerprints and analyzes them. It then returns the result to the device and stores the result on the server for accessing via a web browser or triggering via hooks.

**FAQ (abbreviated)**:
- How does it work? [It uses already available WiFi information to classify locations](https://github.com/schollz/find/blob/master/FAQ.md#how-does-it-work).
- Doesn't this already exist? [Yes, kinda](https://github.com/schollz/find/blob/master/FAQ.md#doesnt-this-already-exist).
- Can I use an iPhone? [Sorry, no](https://github.com/schollz/find/blob/master/FAQ.md#can-i-use-an-iphone).
- Does it work on a Raspberry Pi? [Yes](https://github.com/schollz/find/blob/master/FAQ.md#can-you-run-the-server-on-a-raspberry-pi).
- What's the point? This short piece of code can do [home automation](https://github.com/schollz/find/blob/master/FAQ.md#can-i-run-the-server-at-home-on-the-lan-connection) and [replace motion sensors](https://github.com/schollz/find/blob/master/FAQ.md#can-it-replace-motion-sensors) and  [more](https://github.com/schollz/find/blob/master/FAQ.md#whats-the-point-of-this).
- Can I help? [Yes, please](https://github.com/schollz/find/issues).
- How do I try it? It's easy. Just [download and run](https://github.com/schollz/find#usage).

More questions? See the [unabbreviated FAQ](https://github.com/schollz/find/blob/master/FAQ.md).

# Features

- SSL support
- Compression to keep DBs small
- [Fast](https://github.com/schollz/find/blob/master/BENCHMARKS.md) (20-200x faster than the [previous Python version](https://github.com/schollz/find/tree/python3))
- Mixes two machine learning algorithms for best classifications
- Bug free (yeah, um...probably not. Please [submit an issue](https://github.com/schollz/find/issues) when you find one).

# Requirements

To use this system you need to have the following:
- (Optional) Server: A computer (OS X/Windows/Linux) to run the server. If you don't have this, use [ours](https://ml.internalpositioning.com).
- Client(s): device (laptop/Raspberry Pi/Android smartphone) that has access to WiFi


# Setup (optional)

The tools are prebuilt, so you can skip to the [Usage section](https://github.com/schollz/find#usage) if you just want to try it out.

#### Server
First [install Go](https://golang.org/dl/) if you haven't already. FIND is tested on Go version 1.5+.

```
$ git clone https://github.com/schollz/find.git
$ cd find
$ go get ./...
$ go build
```

Then to run,
```
$ ./find
-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----
   _________  _____
  / __/  _/ |/ / _ \  ______ _____  _____ ____
 / _/_/ //    / // / (_-< -_) __/ |/ / -_) __/
/_/ /___/_/|_/____/ /___|__/_/  |___/\__/_/

(version 2.X) is up and running on http://192.168.1.2:8003
-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----

```

#### Client

The client gathers WiFi fingerprints and sends them to the server. There are two clients - one for an Android smartphone, and one for a WiFi-enabled computer. Check out the individual repos to get started on either/both:
- [FIND app client](https://github.com/schollz/find/tree/android)
- [FIND program client](https://github.com/schollz/find/tree/android)

# Usage

## Gather fingerprint data

To get started using **FIND** you will need one of the client programs. The client programs gather WiFi fingerprints and locations and sends them to a server. There is a client for Android smartphones, and one for WiFi enabled computers.

### Client (for Android smartphones)

*[Sorry iPhone users](https://github.com/schollz/find/blob/master/FAQ.md#can-i-use-an-iphone) but the Apple store prevents apps that access WiFi information, so I will be unable to release a iPhone version.*

To get started using FIND on a smartphone, [download the latest app](https://play.google.com/store/apps/details?id=com.hcp.find) or [build it yourself](https://github.com/schollz/find/tree/android).

When you start up the app you will be asked for a username (enter whatever you want) and you'll be assigned a unique group name. Then you'll see the following:


<center>
<img src="https://i.imgur.com/fbcYom5.png" width="200px" />
<img src="https://i.imgur.com/Ab9eXIk.png" width="200px" />
</center>

Simply click "Learn" and you'll be prompted for a location name. After you enter a location, the app will connect to the server and then submit fingerprints. After you've learned some locations, just hit "Track" and you'll see your calculated location.

<center>
<img src="https://i.imgur.com/fxNIZyP.png" width="200px" />
<img src="https://i.imgur.com/TrgPXns.png" width="200px" />
</center>

To see more detailed information, go to [the server](https://ml.internalpositioning.com) and login with your group name


### Client (for computers)

*Supports Linux, Windows, Raspberry Pi, and OS X!*

To get started, download [the program here](https://github.com/schollz/find/releases/tag/v0.3client) or [build it yourself](https://github.com/schollz/find/tree/fingerprint). To start learning locations simply use

```bash
./fingerprint -e
```

and then to track your location use

```bash
./fingerprint
```

There are other options, you can learn more by typing `./fingerprint --help`. When you start learning/tracking, you can see more detailed information by logging into [the server](https://ml.internalpositioning.com) and login with your group name.

## Analyze fingerprint data

The server analyzes and helps you decipher the fingerprint data, if you are interested in that. Once you got your client you can see statistics about your fingerprints by logging on to the server and signing in with your group name. If you are using our server, login to [ml.internalpositioning.com](https://ml.internalpositioning.com) with your Group name. Otherwise, use your local version of the server.


<center>
<img src="https://i.imgur.com/HSGVyDb.png" width="300px" />
<img src="https://i.imgur.com/IAn5Hss.png" width="300px" />
</center>


You can see fingerprints of individual places by clicking on them, and then you can click on mac addresses to see there statistics across rooms.

<center>
<img src="https://i.imgur.com/3l5UPub.png" width="400px" />
</center>
<center>
<img src="https://i.imgur.com/Udi3xrn.png" width="400px" />
</center>


# Contributing

Please, do! Checkout [the latest issues](https://github.com/schollz/find/issues) to see what needs being done, or add your own cool thing.

If you find a bug or need help with something, feel free to contact:

* Email: [zack@hypercubeplatforms.com](mailto:zack@hypercubeplatforms.com)
* Twitter: [@zack_118](https://twitter.com/intent/tweet?screen_name=zack_118)
* Gitter: [Join room](https://gitter.im/schollz/find)
* Github Issues: [Open an issue](https://github.com/schollz/find/issues/new)

# Acknowledgements

<img src="https://i.imgur.com/Ze51DJ6.png" width="180px" /> Funding from [Duke University Colab](https://colab.duke.edu/)

Thanks to Rishabh Rajgarhia and [CanvasJS](http://canvasjs.com/) for help implementing a nice graph.

Thanks [arafsheikh](https://github.com/arafsheikh) for adding interface selection, [Pugio](https://github.com/Pugio) and [ScottSWu](https://github.com/ScottSWu) for adding OS X/Windows support for the fingerprint program. Thanks [Thom-x](https://github.com/Thom-x) for the Dockerfile. Thanks [certifiedloud](https://github.com/certifiedloud) for implementing the change to `DELETE` requests.

Thanks to [patorjk](http://patorjk.com/software/taag/) and [asciiworld](http://www.asciiworld.com/) for the ASCII art. Thanks to [Imgur](https://imgur.com/a/yjvci) for [hosting](https://imgur.com/a/3yGjV) images.

# Donate

Like this? Help me keep it alive [by donating $5](https://www.paypal.me/ZackScholl/5.00) to [pay for server costs](http://rpiai.com/donate/).

# License

## FIND

**FIND** is a Framework for Internal Navigation and Discovery.

Copyright (C) 2015-2016 Zack Scholl

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the [GNU Affero General Public License](LICENSE) for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see [GNU Affero General Public License here](https://www.gnu.org/licenses/agpl.html).

## CanvasJS

**FIND** uses [CanvasJS](http://canvasjs.com/). Note that you will have to buy the appropriate CanvasJS License if you use this software for commercial purposes. CanvasJS has the following Dual Licensing Model:

### Commercial License

Commercial use of CanvasJS requires you to purchase a license. Without a commercial license you can use it for evaluation purposes only. Please refer to the following link for further details: http://canvasjs.com/.

### Free for Non-Commercial Use

For non-commercial purposes you can use the software for free under Creative Commons Attribution-NonCommercial 3.0 License.

A credit Link is added to the bottom right of the chart which should be preserved. Refer to the following link for further details on the same: http://creativecommons.org/licenses/by-nc/3.0/deed.en_US.
