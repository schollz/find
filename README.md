**The Framework for Internal Navigation and Discovery** (*FIND*) allows you to use your smartphone or laptop to determine your position within your home or office. You can easily use this system in place of motion sensors as its resoltion will allow your phone to distinguish whether you are in the living room or the kitchen or bedroom etc. Simply put, FIND will allow you to replace tons of motion sensors with a single smartphone. The position information can then be used in a variety of ways including home automation, way-finding, tracking, among a few!

<blockquote>Simply put, FIND will allow you to replace <em>tons</em> of motion sensors with a <em>single</em> smartphone!</blockquote>

The system is built on two main components - a server
and a fingerprinting device. The fingerprinting device (computer or android app) sends the specified data to the machine learning server which stores the fingerprints and analyzes them. It then returns the result to the device and stores the result on the server for accessing via a web browser or triggering via hooks.

More detailed documentation can be found in the [FIND Guide](http://internalpositioning.com/guide/).

# Server setup

First get the latest source code:

    git clone https://github.com/schollz/find.git

Installation is very simple. First install Python 3.4 development
packages and start a virtualenv:

    sudo apt-get update
    sudo apt-get -y upgrade
    sudo apt-get install python3.4-dev
    sudo apt-get install python3-pip
    sudo pip3 install virtualenv

    cd find-ml
    virtualenv -p /usr/bin/python3 venv
    source venv/bin/activate

Now you can run the setup using:

    (venv)$ python setup.py install

After which you will be prompted to enter the `address` and `port` of
your server. If you want to run on a home network run `ifconfig` to
check your `address` (it will be something like `192.168.X.Y` usually).
If you want to use an public address you can also use that. If you are
using a reverse proxy you can also set the `external address`, but if
not, you can just leave that blank.

To run the program simple use:

    (venv)$ uwsgi --http address -w server

where `address` is the address you set above.

# App

To use the system, you will need a fingerprinting device. The easiest thing to do is to use [our app from Google Play](https://play.google.com/store/apps/details?id=com.hcp.find&hl=en). 

# Notes

## Backup/restore database

### Backup

```
sqlite3 find.db .sch > schema
sqlite3 find.db .dump > dump
grep -v -f schema dump > data
```

### Restore

```
sqlite3 find.db < data
```

### Copy to new repository

```
rsync -avrP --files-from essential_files ./ ~/find
```