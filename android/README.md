# Source code for making the appIcon

## Requirements

```
sudo npm install -g cordova
```

## Install 
```
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
