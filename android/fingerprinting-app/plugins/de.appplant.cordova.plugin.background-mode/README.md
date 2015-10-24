
<p align="right">
    <a href="https://github.com/katzer/cordova-plugin-background-mode/tree/example">EXAMPLE :point_right:</a>
</p>

Cordova Background Plug-in
==========================

[Cordova][cordova] plugin to prevent the app from going to sleep while in background.

Most mobile operating systems are multitasking capable, but most apps dont need to run while in background and not present for the user. Therefore they pause the app in background mode and resume the app before switching to foreground mode.
The system keeps all network connections open while in background, but does not deliver the data until the app resumes.

### Plugin's Purpose
This cordova plug-in can be used for applications, who rely on continuous network communication independent of from direct user interactions and remote push notifications.

### :bangbang: Store Compliance :bangbang:
The plugin focuses on enterprise-only distribution and may not compliant with all public store vendors.


## Overview
1. [Supported Platforms](#supported-platforms)
2. [Installation](#installation)
3. [ChangeLog](#changelog)
4. [Usage](#usage)
5. [Examples](#examples)
6. [Platform specifics](#platform-specifics)


## Supported Platforms
- __iOS__ (_including iOS8_)
- __Android__ _(SDK >=11)_
- __WP8__


## Installation
The plugin can either be installed from git repository, from local file system through the [Command-line Interface][CLI]. Or cloud based through [PhoneGap Build][PGB].

### Local development environment
From master:
```bash
# ~~ from master branch ~~
cordova plugin add https://github.com/katzer/cordova-plugin-background-mode.git
```
from a local folder:
```bash
# ~~ local folder ~~
cordova plugin add de.appplant.cordova.plugin.background-mode --searchpath path
```
or to use the last stable version:
```bash
# ~~ stable version ~~
cordova plugin add de.appplant.cordova.plugin.background-mode@0.6.3
```

To remove the plug-in, run the following command:
```bash
cordova plugin rm de.appplant.cordova.plugin.background-mode
```

### PhoneGap Build
Add the following xml to your config.xml to always use the latest version of this plugin:
```xml
<gap:plugin name="de.appplant.cordova.plugin.background-mode" version="0.6.3" />
```

More informations can be found [here][PGB_plugin].


## ChangeLog
#### Version 0.6.4 (03.03.2015)
- Resolve possibly dependency conflict

#### Version 0.6.3 (01.01.2015)
- [feature:] Silent mode for Android

#### Version 0.6.2 (14.12.2014)
- [bugfix:] Type error
- [bugfix:] Wrong default values for `isEnabled` and `isActive`.

#### Further informations
- The former `plugin.backgroundMode` namespace has been deprecated and will be removed with the next major release.
- See [CHANGELOG.md][changelog] to get the full changelog for the plugin.

#### Known issues
- Plug-in is broken on Windows Phone 8.1 platform.


## Usage
The plugin creates the object `cordova.plugins.backgroundMode` with  the following methods:

1. [backgroundMode.enable][enable]
2. [backgroundMode.disable][disable]
3. [backgroundMode.isEnabled][is_enabled]
4. [backgroundMode.isActive][is_active]
5. [backgroundMode.getDefaults][android_specifics]
6. [backgroundMode.setDefaults][android_specifics]
7. [backgroundMode.configure][configure]
8. [backgroundMode.onactivate][onactivate]
9. [backgroundMode.ondeactivate][ondeactivate]
10. [backgroundMode.onfailure][onfailure]

### Plugin initialization
The plugin and its methods are not available before the *deviceready* event has been fired.

```javascript
document.addEventListener('deviceready', function () {
    // cordova.plugins.backgroundMode is now available
}, false);
```

### Prevent the app from going to sleep in background
To prevent the app from being paused while in background, the `backroundMode.enable` interface has to be called.

#### Further informations
- The background mode will be activated once the app has entered the background and will be deactivated after the app has entered the foreground.
- To activate the background mode the app needs to be in foreground.

```javascript
cordova.plugins.backgroundMode.enable();
```

### Pause the app while in background
The background mode can be disabled through the `backgroundMode.disable` interface.

#### Further informations
- Once the background mode has been disabled, the app will be paused when in background.

```javascript
cordova.plugins.backgroundMode.disable();
```

### Receive if the background mode is enabled
The `backgroundMode.isEnabled` interface can be used to get the information if the background mode is enabled or disabled.

```javascript
cordova.plugins.backgroundMode.isEnabled(); // => boolean
```

### Receive if the background mode is active
The `backgroundMode.isActive` interface can be used to get the information if the background mode is active.

```javascript
cordova.plugins.backgroundMode.isActive(); // => boolean
```

### Get informed when the background mode has been activated
The `backgroundMode.onactivate` interface can be used to get notified when the background mode has been activated.

```javascript
cordova.plugins.backgroundMode.onactivate = function() {};
```

### Get informed when the background mode has been deactivated
The `backgroundMode.ondeactivate` interface can be used to get notified when the background mode has been deactivated.

#### Further informations
- Once the mode has been deactivated the app will be paused soon after the callback has been fired.

```javascript
cordova.plugins.backgroundMode.ondeactivate = function() {};
```

### Get informed when the background mode could not been activated
The `backgroundMode.onfailure` interface can be used to get notified when the background mode could not been activated.

The listener has to be a function and takes the following arguments:
 - errorCode: Error code which describes the error

```javascript
cordova.plugins.backgroundMode.onfailure = function(errorCode) {};
```


## Examples
The following example demonstrates how to enable the background mode after device is ready. The mode itself will be activated when the app has entered the background.

```javascript
document.addEventListener('deviceready', function () {
    // Android customization
    cordova.plugins.backgroundMode.setDefaults({ text:'Doing heavy tasks.'});
    // Enable background mode
    cordova.plugins.backgroundMode.enable();

    // Called when background mode has been activated
    cordova.plugins.backgroundMode.onactivate = function () {
        setTimeout(function () {
            // Modify the currently displayed notification
            cordova.plugins.backgroundMode.configure({
                text:'Running in background for more than 5s now.'
            });
        }, 5000);
    }
}, false);
```


## Platform specifics

### Android customization
To indicate that the app is executing tasks in background and being paused would disrupt the user, the plug-in has to create a notification while in background - like a download progress bar.

#### Override defaults
The title, ticker and text for that notification can be customized as follows:

```javascript
cordova.plugins.backgroundMode.setDefaults({
    title:  String,
    ticker: String,
    text:   String
})
```

By default the app will come to foreground when taping on the notification. That can be changed also.

```javascript
cordova.plugins.backgroundMode.setDefaults({
    resume: false
})
```

#### Modify the currently displayed notification
It's also possible to modify the currently displayed notification while in background.

```javascript
cordova.plugins.backgroundMode.configure({
    title: String,
    ...
})
```

#### Run in background without notification
In silent mode the plugin will not display a notification - which is not the default. Be aware that Android recommends adding a notification otherwise the OS may pause the app.

```javascript
cordova.plugins.backgroundMode.configure({
    silent: true
})
```


## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request


## License

This software is released under the [Apache 2.0 License][apache2_license].

© 2013-2014 appPlant UG, Inc. All rights reserved


[cordova]: https://cordova.apache.org
[CLI]: http://cordova.apache.org/docs/en/edge/guide_cli_index.md.html#The%20Command-line%20Interface
[PGB]: http://docs.build.phonegap.com/en_US/index.html
[PGB_plugin]: https://build.phonegap.com/plugins/2056
[changelog]: CHANGELOG.md
[enable]: #prevent-the-app-from-going-to-sleep-in-background
[disable]: #pause-the-app-while-in-background
[is_enabled]: #receive-if-the-background-mode-is-enabled
[is_active]: #receive-if-the-background-mode-is-active
[android_specifics]: #android-customization
[configure]: #modify-the-currently-displayed-notification
[onactivate]: #get-informed-when-the-background-mode-has-been-activated
[ondeactivate]: #get-informed-when-the-background-mode-has-been-deactivated
[onfailure]: #get-informed-when-the-background-mode-could-not-been-activated
[apache2_license]: http://opensource.org/licenses/Apache-2.0
[appplant]: http://appplant.de
