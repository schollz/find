cordova-wifiscanner
===================

Cordova interface for [android.net.wifi.WifiManager](http://developer.android.com/reference/android/net/wifi/WifiManager.html) to list available access points and their signal strengths.

This plugin is based on the structure of the [Device-Motion](http://plugins.cordova.io/#/package/org.apache.cordova.device-motion) plugin and exposes a similar API.

Installation
------------

    cordova plugin add nl.nielsad.cordova.wifiscanner

Supported Platforms
-------------------

- Android

Methods
-------

- navigator.wifi.getAccessPoints(successCallback, errorCallback)
- navigator.wifi.watchAccessPoints(successCallback, errorCallback, options)
- navigator.wifi.clearWatch(id)

For more information, see the documentation of the [Device-Motion](http://plugins.cordova.io/#/package/org.apache.cordova.device-motion) plugin.

Objects
-------

- AccessPoint (see [android.net.wifi.ScanResult](http://developer.android.com/reference/android/net/wifi/ScanResult.html))
  - __BSSID__: Address of the access point __(String)__
  - __SSID__: Network name; __(String)__
  - __level__: RSSI __(Number)__
