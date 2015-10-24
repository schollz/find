/*
    Copyright 2013-2014 appPlant UG

    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
*/

var exec    = require('cordova/exec'),
    channel = require('cordova/channel');


// Override back button action to prevent being killed
document.addEventListener('backbutton', function () {}, false);

// Called before 'deviceready' listener will be called
channel.onCordovaReady.subscribe(function () {
    // Device plugin is ready now
    channel.onCordovaInfoReady.subscribe(function () {
        // Set the defaults
        exports.setDefaults({});
    });

    // Only enable WP8 by default
    if (['WinCE', 'Win32NT'].indexOf(device.platform) > -1) {
        exports.enable();
    }
});


/**
 * @private
 *
 * Flag indicated if the mode is enabled.
 */
exports._isEnabled = false;

/**
 * @private
 *
 * Flag indicated if the mode is active.
 */
exports._isActive = false;

/**
 * @private
 *
 * Default values of all available options.
 */
exports._defaults = {
    title:  'App is running in background',
    text:   'Doing heavy tasks.',
    ticker: 'App is running in background',
    resume: true,
    silent: false
};


/**
 * Activates the background mode. When activated the application
 * will be prevented from going to sleep while in background
 * for the next time.
 */
exports.enable = function () {
    this._isEnabled = true;
    cordova.exec(null, null, 'BackgroundMode', 'enable', []);
};

/**
 * Deactivates the background mode. When deactivated the application
 * will not stay awake while in background.
 */
exports.disable = function () {
    this._isEnabled = false;
    cordova.exec(null, null, 'BackgroundMode', 'disable', []);
};

/**
 * List of all available options with their default value.
 *
 * @return {Object}
 */
exports.getDefaults = function () {
    return this._defaults;
};

/**
 * Overwrite default settings
 *
 * @param {Object} overrides
 *      Dict of options which shall be overridden
 */
exports.setDefaults = function (overrides) {
    var defaults = this.getDefaults();

    for (var key in defaults) {
        if (overrides.hasOwnProperty(key)) {
            defaults[key] = overrides[key];
        }
    }

    if (device.platform == 'Android') {
        cordova.exec(null, null, 'BackgroundMode', 'configure', [defaults, false]);
    }
};

/**
 * Configures the notification settings for Android.
 * Will be merged with the defaults.
 *
 * @param {Object} options
 *      Dict with key/value pairs
 */
exports.configure = function (options) {
    var settings = this.mergeWithDefaults(options);

    if (device.platform == 'Android') {
        cordova.exec(null, null, 'BackgroundMode', 'configure', [settings, true]);
    }
};

/**
 * If the mode is enabled or disabled.
 *
 * @return {Boolean}
 */
exports.isEnabled = function () {
    return this._isEnabled;
};

/**
 * If the mode is active.
 *
 * @return {Boolean}
 */
exports.isActive = function () {
    return this._isActive;
};

/**
 * Called when the background mode has been activated.
 */
exports.onactivate = function () {};

/**
 * Called when the background mode has been deaktivated.
 */
exports.ondeactivate = function () {};

/**
 * Called when the background mode could not been activated.
 *
 * @param {Integer} errorCode
 *      Error code which describes the error
 */
exports.onfailure = function () {};

/**
 * @private
 *
 * Merge settings with default values.
 *
 * @param {Object} options
 *      The custom options
 *
 * @return {Object}
 *      Default values merged
 *      with custom values
 */
exports.mergeWithDefaults = function (options) {
    var defaults = this.getDefaults();

    for (var key in defaults) {
        if (!options.hasOwnProperty(key)) {
            options[key] = defaults[key];
            continue;
        }
    }

    return options;
};
