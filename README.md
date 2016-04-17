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


cp ../appIcon.png platforms/android/res/drawable-hdpi/icon.png
cp ../appIcon.png platforms/android/res/drawable-mdpi/icon.png
cp ../appIcon.png platforms/android/res/drawable-ldpi/icon.png
cp ../appIcon.png platforms/android/res/drawable-xhdpi/icon.png

cp ../index.html platforms/android/assets/www/ && cp ../jquery-1.9.js platforms/android/assets/www/ && cp ../main.js platforms/android/assets/www/

./platforms/android/cordova/run --device

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

Copyright 2015-2016 Zack Scholl

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License [https://github.com/schollz/find/blob/master/LICENSE](https://github.com/schollz/find/blob/master/LICENSE).

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
