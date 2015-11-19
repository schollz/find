**The Framework for Internal Navigation and Discovery** (*FIND*) allows you to use your smartphone or laptop to determine your position within your home or office. You can easily use this system in place of motion sensors as its resoltion will allow your phone to distinguish whether you are in the living room or the kitchen or bedroom etc. Simply put, FIND will allow you to replace tons of motion sensors with a single smartphone. The position information can then be used in a variety of ways including home automation, way-finding, tracking, among a few!

<blockquote>Simply put, FIND will allow you to replace <em>tons</em> of motion sensors with a <em>single</em> smartphone!</blockquote>

The system is built on two main components - a server
and a fingerprinting device. The fingerprinting device (computer or android app) sends the specified data to the machine learning server which stores the fingerprints and analyzes them. It then returns the result to the device and stores the result on the server for accessing via a web browser or triggering via hooks.

More detailed documentation can be found in the [FIND Guide](http://internalpositioning.com/guide/).

# Requirements

To use this system you need to have the following

- Linux / Mac / Cygwin (Windows). Windows is not yet supported (but will be soon). Raspberry Pi is supported!
- Python3 installed
- Either WiFi capable laptop or Android smartphone

# Setup

Installation is very simple. Simply download the latest source code and install:

    git clone https://github.com/schollz/find.git
    cd find/
    sudo python3 setup.py 

After which you will be prompted to enter the `address` and `port` of
your server. If you want to run on a home network run `ifconfig` to
check your address (it will be something like `192.168.X.Y` usually).
If you want to use an public address you can also use that. Note: if you are
using a reverse proxy you can also set the `external address`, but if
not, you can just leave that blank.

To run **FIND** simply use:

    python3 server.py

To actually use the system you will need a fingerprinting device. The easiest thing to do is to use [our app from Google Play](https://play.google.com/store/apps/details?id=com.hcp.find&hl=en) or [build the app yourself](http://internalpositioning.com/guide/deploy/#building-android-app). Alternatively, you don't have to build an app at all and can simply use your laptop via the [the fingerprinting program](https://github.com/schollz/find/blob/master/computer/fingerprinting.py), more details [here](http://internalpositioning.com/guide/deploy/#laptop-computer).


# Using FIND

## 1. Fingerprint locations

### If you want to use the app...

First [download the app from the Google Play store](https://play.google.com/store/apps/details?id=com.hcp.find). 

![Guide to starting app](http://www.internalpositioning.com/guide/img/guide_app_guide_1.png)

To begin fingerprinting, stand in a location and enter the location name into the app. Then hit the "learn" button in the app. Then simply press ```ON```. The app will then run at the specified interval, even in the background. To stop the fingerprinting you have to press ```OFF``` and to stop the program entirely you have to press ```EXIT```.


### If you want to use a laptop...

Any computer with a WiFi card or laptops will be able to use FIND as well. Its simple to get started. If you cloned the repository, goto ```computer/``` to find ```fingerprinting.py```. If you haven't cloned the repository, you can simply [download the fingerprinting.py script]([the fingerprinting program](https://github.com/schollz/find/blob/master/computer/fingerprinting.py)).

To fingerprint locations simply use

```bash
python3 fingerprinting.py -l "name of location" -u "user name" -g "your group" -s "server" -p "port" -c 10 -r learn
```

which will send 10 fingerprints to the server. If you are using the demo server, simply set "server" to "finddemo.duckdns.org" and do not include the port. If you are running locally you probably want "server" to be "localhost" and port to be whatever you specified. The name of "your group" can be whatever you want.

Repeat this process for a few locations.


## 2. Analyze fingerprints

Now that you have learned several fingerprints, open a web browser and
navigate to the dashboard page at `http://address:port/` or [http://finddemo.duckdns.org](http://finddemo.duckdns.org) if you are using the demo server. Login with the name of your group that you specified in the app or in the script.

Once you login you'll be able to access the "Dashboard." This dashboard page contains all the information about the learned fingerprints and the analysis. More information about the dashboard page can be found on the [API documentation](/api/#webpages).

The dashboard has many options and edits that you can do. For now, the only thing you need to do is press the button ```Calculate All and Find Parameters``` which will automatically optimize the parameters and generate the dataset you need for tracking.

![Guide to analyzing fingerprints with app](http://www.internalpositioning.com/guide/img/guide_dashboard.png)

## 3. Track locations

To see your current position classification, simply hit "Classifications" from the webpage that you visited to see the dashboard. This classifications are automatically updated as new information is available from the app/laptop. Sending the tracking information is very easy:

### If you want to use the app...


Simply go back to the app and click the "track" button and then hit ```ON```. Now you are tracking!

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

![Screenshot of the signin](http://internalpositioning.com/guide/img/signin1.png)

</center>

<br>

<center>

![Landing](http://internalpositioning.com/guide/img/landing2.png)
*Screenshot of the landing page*

</center><br>

## Monitor location in realtime


<center>

![Screenshot of the classifications page](http://internalpositioning.com/guide/img/classifications1.png)

</center><br>

## Visualize accuracy and errors

<center>

![Charts show a clear diagnostics of the accuracy for each room](http://internalpositioning.com/guide/img/stats1.png)

</center><br>

<center>

![Pie charts lets you visualize the classification errors](http://internalpositioning.com/guide/img/pies1.png)

</center><br>

## Visualize raw data


<center>

![In-depth analysis of the raw fingerprint data](http://internalpositioning.com/guide/img/signals1.png)

</center><br>


## Tip of the iceberg

![Tip of the iceberg](http://www.internalpositioning.com/guide/img/iceberg.png)

There are lots of other features available which you can start investigating when you get used to the app and server. Some things to check out:

- [Build your own app](http://www.internalpositioning.com/guide/deploy/#building-android-app) with internal positioning builtin.
- [Use the RESTFUL API](http://www.internalpositioning.com/guide/api/#routes) for your own applications, like home automation.
- [Deploy the server at home](http://www.internalpositioning.com/guide/deploy/#server-setup).
