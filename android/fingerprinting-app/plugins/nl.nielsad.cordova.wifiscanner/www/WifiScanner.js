/**
 * Author: Niels A.D.
 * Project: cordova-wifiscanner (https://github.com/nielsAD/cordova-wifiscanner)
 * License: Apache License v2.0 (http://www.apache.org/licenses/LICENSE-2.0)
 *
 * JavaScript interface to the native scanner
 * Based on the Device-Motion Cordova plugin
 */

/**
 * This class provides access to the available access points.
 * @constructor
 */
var argscheck = require('cordova/argscheck'),
    utils = require("cordova/utils"),
    exec = require("cordova/exec");//,
    AccessPoint = require('./AccessPoint');

// Is the adapter listening?
var running = false;

// Keeps reference to watchAccessPoints calls.
var timers = {};

// Array of listeners; used to keep track of when we should call start and stop.
var listeners = [];

// Last returned network scan from native
var networks = null;

// Tells native to start.
function start() {
    exec(function(res) {
        var tempListeners = listeners.slice(0);
        networks = res.map(function(r){ return new AccessPoint(r.BSSID, r.SSID, r.level); });
        for (var i = 0, l = tempListeners.length; i < l; i++) {
            tempListeners[i].win(networks);
        }
    }, function(e) {
        var tempListeners = listeners.slice(0);
        for (var i = 0, l = tempListeners.length; i < l; i++) {
            tempListeners[i].fail(e);
        }
    }, "WifiScanner", "start", []);
    running = true;
}

// Tells native to stop.
function stop() {
    exec(null, null, "WifiScanner", "stop", []);
    running = false;
}

// Adds a callback pair to the listeners array
function createCallbackPair(win, fail) {
    return {win:win, fail:fail};
}

// Removes a win/fail listener pair from the listeners array
function removeListeners(l) {
    var idx = listeners.indexOf(l);
    if (idx > -1) {
        listeners.splice(idx, 1);
        if (listeners.length === 0) {
            stop();
        }
    }
}

var wifi = {
    /**
     * Asynchronously acquires the current list of access points.
     *
     * @param {Function} successCallback      The function to call when the network scan is available
     * @param {Function} errorCallback        The function to call when there is an error getting the access points. (OPTIONAL)
     * @param {AccelerationOptions} options   The options for getting the wifi data such as timeout. (OPTIONAL)
     */
    getAccessPoints: function(successCallback, errorCallback, options) {
        argscheck.checkArgs('fFO', 'wifiscanner.getAccessPoints', arguments);

        var p;
        var win = function(a) {
            removeListeners(p);
            successCallback(a);
        };
        var fail = function(e) {
            removeListeners(p);
            errorCallback && errorCallback(e);
        };

        p = createCallbackPair(win, fail);
        listeners.push(p);

        if (!running) {
            start();
        }
    },

    /**
     * Asynchronously acquires the acceleration repeatedly at a given interval.
     *
     * @param {Function} successCallback   The function to call each time the network scan is available
     * @param {Function} errorCallback     The function to call when there is an error getting the access points. (OPTIONAL)
     * @param {WifiOptions} options        The options for getting the wifi data such as timeout. (OPTIONAL)
     * @return String                      The watch id that must be passed to #clearWatch to stop watching.
     */
    watchAccessPoints: function(successCallback, errorCallback, options) {
        argscheck.checkArgs('fFO', 'wifiscanner.watchAccessPoints', arguments);
        // Default interval (10 sec)
        var frequency = (options && options.frequency && typeof options.frequency == 'number') ? options.frequency : 10000;

        // Keep reference to watch id, and report networks readings as often as defined in frequency
        var id = utils.createUUID();

        var p = createCallbackPair(function(){}, function(e) {
            removeListeners(p);
            errorCallback && errorCallback(e);
        });
        listeners.push(p);

        timers[id] = {
            timer: window.setInterval(function() {
                if (networks) {
                    successCallback(networks);
                }
            }, frequency),
            listeners: p
        };

        if (running) {
            // If we're already running then immediately invoke the success callback
            // but only if we have retrieved a value, sample code does not check for null ...
            if (networks) {
                successCallback(networks);
            }
        } else {
            start();
        }

        return id;
    },

    /**
     * Clears the specified wifi watch.
     *
     * @param {String} id   The id of the watch returned from #watchAccessPoints.
     */
    clearWatch: function(id) {
        // Stop javascript timer & remove from timer list
        if (id && timers[id]) {
            window.clearInterval(timers[id].timer);
            removeListeners(timers[id].listeners);
            delete timers[id];
        }
    }
};
module.exports = wifi;
