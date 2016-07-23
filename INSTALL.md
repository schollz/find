# Installation


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

If you'd like to add (experimental) SVM support, please [see the SVM docs](https://doc.internalpositioning.com/svm/).

If you'd like to add MQTT support, please [see the MQTT docs](https://doc.internalpositioning.com/mqtt/).
