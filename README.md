# Source code for making the app

You don't need to build this, you can [use ours](https://play.google.com/store/apps/details?id=com.hcp.find).

## Requirements

```
sudo npm install -g cordova
```

## Install
```
git clone https://github.com/schollz/find.git
cd find
git checkout android
cordova create find com.hcp.find Find
cd find
cordova platform add android
cordova plugin add https://github.com/schollz/cordova-plugin-wifi.git
cordova plugin add cordova-plugin-whitelist
cordova plugin add https://github.com/schollz/cordova-plugin-background-mode.git
cordova plugin add cordova-plugin-dialogs
cordova plugin add cordova-plugin-powermanagement
cordova plugin add https://github.com/schollz/phonegap-plugin-brightness.git
cordova plugin add cordova-plugin-chrome-apps-alarms
cordova plugins add https://github.com/uniclau/AlarmPlugin.git
cordova plugin add cordova-plugin-vibration
cp ../index.html platforms/android/assets/www/ && cp ../jquery-1.9.js platforms/android/assets/www/ && cp ../main.js platforms/android/assets/www/ && ./platforms/android/cordova/build

./platforms/android/cordova/run --device

cp ../appIcon.png platforms/android/res/drawable-hdpi/icon.png
cp ../appIcon.png platforms/android/res/drawable-mdpi/icon.png
cp ../appIcon.png platforms/android/res/drawable-ldpi/icon.png
cp ../appIcon.png platforms/android/res/drawable-xhdpi/icon.png



./platforms/android/cordova/build
./platforms/android/cordova/run --device
cp /home/phi/Downloads/cord/find/platforms/android/build/outputs/apk/android-debug.apk ~/Dropbox/android-debug.apk
```


To build actual app:

```
Open Android Studio
Import Project Gradle find->platforms->android
Change version code and version name in android/manifests/AndroidManifest.xml
Build signed APK using keystore Dropbox/keystore/alskasldfk.jks
```

# License

FIND app is the app for the Framework for Internal Navigation and Discovery.

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
