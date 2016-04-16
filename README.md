[![Join the chat at https://gitter.im/schollz/find](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/schollz/find?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) [![Version 0.21prealpha](https://img.shields.io/badge/version-0.21prealpha-brightgreen.svg)](https://www.internalpositioning.com/guide/development/)
<center>

![Image](https://www.internalpositioning.com/guide/img/wifi-marker-darkgrey-small.png
)

</center>

*Note for past users, the codebase has been completely rewritten in Golang so [things are ~100x faster](https://github.com/schollz/find/blob/master/BENCHMARKS.md#python-vs-go), smarter, and more secure. The [Python3 version will stay available](https://github.com/schollz/find/tree/python3), but it is no longer supported. Development will now be focused on the [Golang version](https://github.com/schollz/find/tree/master).*

- [Requirements](#requirements)
- [Server setup](#server-setup)
- [Using FIND](#using-find)
  - [1. Fingerprint locations](#1-fingerprint-locations)
  - [2. Analyze fingerprints](#2-analyze-fingerprints)
  - [3. Track locations](#3-track-locations)

- [Screenshots](#screenshots)
- [Contact us](#contact-us)

**The Framework for Internal Navigation and Discovery** (_FIND_) allows you to use your smartphone or laptop to determine your position within your home or office. You can easily use this system in place of motion sensors as its resolution will allow your phone to distinguish whether you are in the living room, the kitchen or the bedroom, etc. The position information can then be used in a variety of ways including home automation, way-finding, or tracking!
<blockquote>Simply put, FIND will allow you to replace <em>tons</em> of motion sensors with a <em>single</em> smartphone!</blockquote>

The system is built on two main components - a server and a fingerprinting device. The fingerprinting device (computer or android app) sends the specified data to the machine learning server which stores the fingerprints and analyzes them. It then returns the result to the device and stores the result on the server for accessing via a web browser or triggering via hooks.

More information? Check out [our splash page overview](http://www.internalpositioning.com/), the [FAQ](https://www.internalpositioning.com/guide/faq/), and the [API](http://www.internalpositioning.com/guide/api/).

# Requirements
To use this system you need to have the following
- Linux / Mac / Cygwin (Windows). Windows is not yet supported (but will be soon). Raspberry Pi is supported!
- Python3 installed
- Either WiFi capable laptop or Android smartphone
- (Optional) Android Studio installed ([in case you want to build the app yourself](https://www.internalpositioning.com/guide/deploy/#building-android-app))

# Server setup
_Note: you don't have to setup a server at all. If you'd like, you can use [our demo server](http://finddemo.duckdns.org) - simply [follow the quickstart guide](https://www.internalpositioning.com/guide/getting-started/) to get going._

Installation is very simple. Simply download the latest source code:

```
git clone https://github.com/schollz/find.git
cd find/
sudo python3 setup.py
```

and then install:

_Note: when you run the installation you will be downloading binaries containing the classification (stored in [releases](https://github.com/schollz/find/releases)). This is the only part of the code that is not open. The classification algorithm is compiled in an attempt to obfuscate how it works. Even though the algorithm is solely and unequivocally my own original creation, there are currently tons of patents being submitted about Wifi-based positioning. I don't want to be responsible for possibly causing people to accidentally infringe on patent-holders, so I opted to make it more difficult for a patent holder to find my exact implementation. I suppose I could check all the patents to make sure I'm not infringing, but I rather stay ignorant to what they've done in order to better my case that mine is an original work._

```
sudo python3 setup.py
```

After which you will be prompted to enter the `address` and `port` of your server. If you want to run on a home network run `ifconfig` to check your address (it will be something like `192.168.X.Y` usually). If you want to use an public address you can also use that. Note: if you are using a reverse proxy you can also set the `external address`, but if not, you can just leave that blank.

To run **FIND** simply use:

```
python3 server.py
```

To actually use the system you will need a fingerprinting device. The easiest thing to do is to use [our app from Google Play](https://play.google.com/store/apps/details?id=com.hcp.find&hl=en) or [build the app yourself](https://www.internalpositioning.com/guide/deploy/#building-android-app). Alternatively, you don't have to build an app at all and can simply use your laptop via the [the fingerprinting program](https://github.com/schollz/find/blob/master/computer/fingerprinting.py), more details [here](https://www.internalpositioning.com/guide/deploy/#laptop-computer).

# Using FIND
## 1. Fingerprint locations
### If you want to use the app...
First [download the app from the Google Play store](https://play.google.com/store/apps/details?id=com.hcp.find).

![Guide to starting app](http://www.internalpositioning.com/guide/img/guide_app_guide_1.png)

To begin fingerprinting, stand in a location and enter the location name into the app. Then hit the "learn" button in the app. Then simply press `ON`. The app will then run at the specified interval, even in the background. To stop the fingerprinting you have to press `OFF` and to stop the program entirely you have to press `EXIT`.

### If you want to use a laptop...
Any computer with a WiFi card or laptops will be able to use FIND as well. Its simple to get started. If you cloned the repository, goto `computer/` to find `fingerprinting.py`. If you haven't cloned the repository, you can simply [download the fingerprinting.py script]([the fingerprinting program](https://github.com/schollz/find/blob/master/computer/fingerprinting.py)).

To fingerprint locations simply use

```bash
python3 fingerprinting.py -l "some location" -u "some user" -g "some group" -s "localhost" -p 8888 -c 10 -r learn
```

which will send 10 fingerprints to the server located at `localhost:8888` of "some location" for "some user" in "some group." If you are using the demo server, simply set "server" to `finddemo.duckdns.org` and do not include the port. If you are running locally you probably want "server" to be "localhost" and port to be whatever you specified. The name of "your group" can be whatever you want.

Repeat this process for a few different locations, making sure to change "some location" to whatever location you are currently located.

## 2. Analyze fingerprints
Now that you have learned several fingerprints, open a web browser and navigate to the dashboard page at `http://address:port/` or [http://finddemo.duckdns.org](http://finddemo.duckdns.org) if you are using the demo server. Login with the name of your group that you specified in the app or in the script.

Once you login you'll be able to access the "Dashboard." This dashboard page contains all the information about the learned fingerprints and the analysis. More information about the dashboard page can be found on the [API documentation](/api/#webpages).

The dashboard has many options and edits that you can do. For now, the only thing you need to do is press the button `Calculate All and Find Parameters` which will automatically optimize the parameters and generate the dataset you need for tracking.

![Guide to analyzing fingerprints with app](http://www.internalpositioning.com/guide/img/guide_dashboard.png)

## 3. Track locations
To see your current position classification, simply hit "Classifications" from the webpage that you visited to see the dashboard. This classifications are automatically updated as new information is available from the app/laptop. Sending the tracking information is very easy:

### If you want to use the app...
Simply go back to the app and click the "track" button and then hit `ON`. Now you are tracking!

![Guide to analyzing fingerprints](http://www.internalpositioning.com/guide/img/guide_tracking.png)

### If you want to use a laptop...
To track locations simply use

```bash
python3 fingerprinting.py -u "user name" -g "your group" -s "server" -p "port" -c 1000 -r track
```

which will run 1000 times.

# Screenshots
## Sign-in
<center>

![Screenshot of the signin](https://www.internalpositioning.com/guide/img/signin1.png)

</center>

<br>
<center>

![Landing](https://www.internalpositioning.com/guide/img/landing2.png) _Screenshot of the landing page_

</center><br>

## Monitor location in realtime
<center>

![Screenshot of the classifications page](https://www.internalpositioning.com/guide/img/classifications1.png)

</center><br>

## Visualize accuracy and errors
<center>

![Charts show a clear diagnostics of the accuracy for each room](https://www.internalpositioning.com/guide/img/stats1.png)

</center><br>
<center>

![Pie charts lets you visualize the classification errors](https://www.internalpositioning.com/guide/img/pies1.png)

</center><br>

## Visualize raw data
<center>

![In-depth analysis of the raw fingerprint data](https://www.internalpositioning.com/guide/img/signals1.png)

</center><br>

## Tip of the iceberg
![Tip of the iceberg](http://www.internalpositioning.com/guide/img/iceberg.png)

There are lots of other features available which you can start investigating when you get used to the app and server. Some things to check out:
- [Build your own app](http://www.internalpositioning.com/guide/deploy/#building-android-app) with internal positioning builtin.
- [Use the RESTFUL API](http://www.internalpositioning.com/guide/api/#routes) for your own applications, like home automation.
- [Deploy the server at home](http://www.internalpositioning.com/guide/deploy/#server-setup).

# Contact us
- Email: [zack@hypercubeplatforms.com](zack@hypercubeplatforms.com)
- Gitter: [Join room](https://gitter.im/schollz/find)
- Github Issues: [Open an issue](https://github.com/schollz/find/issues/new)
- Subscribe: [Get latest updates](http://hypercubeplatforms.us10.list-manage1.com/subscribe?u=885d1826479b36238603d935c&id=dfc8e534c4)
