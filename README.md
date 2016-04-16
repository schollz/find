# ![Image](https://www.internalpositioning.com/guide/img/wifi-marker-darkgrey-small.png) FIND

[![Join the chat at https://gitter.im/schollz/find](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/schollz/find?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) [![Version 2.0](https://img.shields.io/badge/version-2.0-brightgreen.svg)](https://www.internalpositioning.com/guide/development/)
[![Go Report Card](https://goreportcard.com/badge/github.com/schollz/gofind)](https://goreportcard.com/report/github.com/schollz/gofind)


**The Framework for Internal Navigation and Discovery** (_FIND_) allows you to use your (Android) smartphone or laptop to determine your position within your home or office. You can easily use this system in place of motion sensors as its resolution will allow your phone to distinguish whether you are in the living room, the kitchen or the bedroom, etc. The position information can then be used in a variety of ways including home automation, way-finding, or tracking!
<blockquote>Simply put, FIND will allow you to replace <em>tons</em> of motion sensors with a <em>single</em> smartphone!</blockquote>

The system is built on two main components - a server and a fingerprinting device. The fingerprinting device (computer or android app) sends the specified data to the machine learning server which stores the fingerprints and analyzes them. It then returns the result to the device and stores the result on the server for accessing via a web browser or triggering via hooks.


# Features

- SSL support
- Compression to keep DBs small
- [Fast]() (20-200x faster than the [previous Python version]())

# Requirements
To use this system you need to have the following:
- A device (laptop/Raspberry Pi/Android smartphone) that has access to WiFi
- (Optional) A computer (OS X/Windows/Linux) to run the server. If you don't have this, use [ours](https://ml.internalpositioning.com).


# Setup

## 1. Server (optional)

_Note:_ You are welcome to skip this step and use [our server](https://ml.internalpositioning.com), just make sure to point the server address to https://ml.internalpositioning.com.

```bash
$ git clone https://github.com/schollz/gofind.git
$ cd gofind
$ go get ./...
$ go build
$ ./gofind
-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----
               ,_   .  ._. _.  .
           , _-\','|~\~      ~/      ;-'_   _-'     ,;_;_,    ~~-
  /~~-\_/-'~'--' \~~| ',    ,'      /  / ~|-_\_/~/~      ~~--~~~~'--_
  /              ,/'-/~ '\ ,' _  , '|,'|~                   ._/-, /~
  ~/-'~\_,       '-,| '|. '   ~  ,\ /'~                /    /_  /~
.-~      '|        '',\~|\       _\~     ,_  ,               /|
          '\        /'~          |_/~\\,-,~  \ "         ,_,/ |
           |       /            ._-~'\_ _~|              \ ) /
            \   __-\           '/      ~ |\  \_          /  ~
  .,         '\ |,  ~-_      - |          \\_' ~|  /\  \~ ,
               ~-_'  _;       '\           '-,   \,' /\/  |
                 '\_,~'\_       \_ _,       /'    '  |, /|'
                   /     \_       ~ |      /         \  ~'; -,_.
                   |       ~\        |    |  ,        '-_, ,; ~ ~\
                    \,      /        \    / /|            ,-, ,   -,
                     |    ,/          |  |' |/          ,-   ~ \   '.
                    ,|   ,/           \ ,/              \       |
                    /    |             ~                 -~~-, /   _
                    |  ,-'                                    ~    /
                    / ,'                                      ~
                    ',|  ~
                      ~'
-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----
   _________  _____
  / __/  _/ |/ / _ \  ______ _____  _____ ____
 / _/_/ //    / // / (_-< -_) __/ |/ / -_) __/
/_/ /___/_/|_/____/ /___|__/_/  |___/\__/_/

(version 2.X) is up and running on http://192.168.1.2:8003
-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----

```

## 2. Client

The client gathers WiFi fingerprints and sends them to the server.

If you want to use an Android device,  [download our app](https://play.google.com/store/apps/details?id=com.hcp.find), or [build the app yourself]().

If you want to use a computer (laptop/Raspberry Pi/etc.), [download our client](), or [build it yourself]().

# Usage

## App

When you start up the app you will be asked for a username (enter whatever you want) and you'll be assigned a unique group name. Then you'll see the following:

![App1](https://i.imgur.com/bGVVQeW.png =100x)

Simply click "Learn" and you'll be prompted for a location name.

![App2](https://i.imgur.com/cqab0bl.png)

After you enter a location, the app will connect to the server and then submit fingerprints.

![App3](https://i.imgur.com/kwwLVGL.png)

After you've learned some locations, just hit "Track" and you'll see your calculated location.

![App4](https://i.imgur.com/3mMV7e7.png)

To see more detailed information, goto [the server](https://ml.internalpositioning.com) and login with your group name


## Client program

To start learning locations simply use

```bash
./fingerprint -e
```

and then to track your location use

```bash
./fingerprint
```

There are other options, you can learn more by [looking at the README](). To see more detailed information, goto [the server](https://ml.internalpositioning.com) and login with your group name

## Server

You can see statistics about your fingeprints by logging on to the server and signing in with your group name.

![stats](https://i.imgur.com/HSGVyDb.jpg)

You can see fingerprints of individual places by clicking on them

![places](https://i.imgur.com/3l5UPub.jpg)

and then you can click on mac addresses to see there statistics across rooms

![macs](https://i.imgur.com/Udi3xrn.jpg)

The server also reveals realtime tracking

![realtime](https://i.imgur.com/IAn5Hss.jpg)


# Acknowledgements

Thanks to [patorjk](http://patorjk.com/software/taag/) and [asciiworld](http://www.asciiworld.com/) for the ASCII art.

[Imgur](https://imgur.com/a/yjvci) for [hosting](https://imgur.com/a/3yGjV) images
