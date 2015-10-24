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

using WPCordovaClassLib.Cordova.Commands;
using Windows.Devices.Geolocation;
using Microsoft.Phone.Shell;
using System;
using WPCordovaClassLib.Cordova;

namespace Cordova.Extension.Commands
{
    /// </summary>
    /// Ermöglicht, dass eine Anwendung im Hintergrund läuft ohne pausiert zu werden
    /// </summary>
    public class BackgroundMode : BaseCommand
    {
        /// </summary>
        /// Event types for callbacks
        /// </summary>
        enum Event {
            ACTIVATE, DEACTIVATE, FAILURE
        }

        #region Instance variables

        /// </summary>
        /// Flag indicates if the plugin is enabled or disabled
        /// </summary>
        private bool IsDisabled = true;

        /// </summary>
        /// Geolocator to monitor location changes
        /// </summary>
        private static Geolocator Geolocator { get; set; }

        #endregion

        #region Interface methods

        /// </summary>
        /// Enable the mode to stay awake when switching
        /// to background for the next time.
        /// </summary>
        public void enable (string args)
        {
            IsDisabled = false;
        }

        /// </summary>
        /// Disable the background mode and stop
        /// being active in background.
        /// </summary>
        public void disable (string args)
        {
            IsDisabled = true;

            Deactivate();
        }

        #endregion

        #region Core methods

        /// </summary>
        /// Keep the app awake by tracking
        /// for position changes.
        /// </summary>
        private void Activate()
        {
            if (IsDisabled || Geolocator != null)
                return;

            if (!IsServiceAvailable())
            {
                FireEvent(Event.FAILURE, null);
                return;
            }

            Geolocator = new Geolocator();

            Geolocator.DesiredAccuracy   = PositionAccuracy.Default;
            Geolocator.MovementThreshold = 100000;
            Geolocator.PositionChanged  += geolocator_PositionChanged;

            FireEvent(Event.ACTIVATE, null);
        }

        /// </summary>
        /// Let the app going to sleep.
        /// </summary>
        private void Deactivate ()
        {
            if (Geolocator == null)
                return;

            FireEvent(Event.DEACTIVATE, null);

            Geolocator.PositionChanged -= geolocator_PositionChanged;
            Geolocator = null;
        }

        #endregion

        #region Helper methods

        /// </summary>
        /// Determine if location service is available and enabled.
        /// </summary>
        private bool IsServiceAvailable()
        {
            Geolocator geolocator = (Geolocator == null) ? new Geolocator() : Geolocator;

            PositionStatus status = geolocator.LocationStatus;

            if (status == PositionStatus.Disabled)
                return false;

            if (status == PositionStatus.NotAvailable)
                return false;

            return true;
        }

        /// <summary>
        /// Fires the given event.
        /// </summary>
        private void FireEvent(Event Event, string Param)
        {
            string EventName;

            switch (Event) {
                case Event.ACTIVATE:
                    EventName = "activate"; break;
                case Event.DEACTIVATE:
                    EventName = "deactivate"; break;
                default:
                    EventName = "failure"; break;
            }

            string js = String.Format("cordova.plugins.backgroundMode.on{0}({1})", EventName, Param);

            PluginResult pluginResult = new PluginResult(PluginResult.Status.OK, js);

            pluginResult.KeepCallback = true;

            DispatchCommandResult(pluginResult);
        }

        #endregion

        #region Delegate methods

        private void geolocator_PositionChanged(Geolocator sender, PositionChangedEventArgs args)
        {
            // Nothing to do here
        }

        #endregion

        #region Lifecycle methods

        /// <summary>
        /// Occurs when the application is being deactivated.
        /// </summary>
        public override void OnPause(object sender, DeactivatedEventArgs e)
        {
            Activate();
        }

        /// <summary>
        /// Occurs when the application is being made active after previously being put
        /// into a dormant state or tombstoned.
        /// </summary>
        public override void OnResume(object sender, ActivatedEventArgs e)
        {
            Deactivate();
        }

        #endregion
    }
}
