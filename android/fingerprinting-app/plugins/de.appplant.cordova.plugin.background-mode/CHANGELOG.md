## ChangeLog
#### Version 0.6.4 (03.03.2015)
- Resolve possibly dependency conflict

#### Version 0.6.3 (01.01.2015)
- [feature:] Silent mode for Android

#### Version 0.6.2 (14.12.2014)
- [bugfix:] Type error
- [bugfix:] Wrong default values for `isEnabled` and `isActive`.

#### Version 0.6.1 (14.12.2014)
- [enhancement:] Set default settings through `setDefaults`.
- [enhancement:] New method `isEnabled` to receive if mode is enabled.
- [enhancement:] New method `isActive` to receive if mode is active.
- [bugfix:] Events caused thread collision.


#### Version 0.6.0 (14.12.2014)
- [feature:] Android support
- [feature:] Change Android notification through `configure`.
- [feature:] `onactivate`, `ondeactivate` and `onfailure` callbacks.
- [___change___:] Disabled by default
- [enhancement:] Get default settings through `getDefaults`.
- [enhancement:] iOS does not require user permissions, internet connection and geo location anymore.

#### Version 0.5.0 (13.02.2014)
- __retired__

#### Version 0.4.1 (13.02.2014)
- Release under the Apache 2.0 license.
- [enhancement:] Location tracking is only activated on WP8 if the location service is available.
- [bigfix:] Nullpointer exception on WP8.

#### Version 0.4.0 (10.10.2013)
- Added WP8 support<br>
  The plugin turns the app into an location tracking app *(for the time it runs in the background)*.

#### Version 0.2.1 (09.10.2013)
- Added js interface to manually enable/disable the background mode.

#### Version 0.2.0 (08.10.2013)
- Added iOS (>= 5) support<br>
  The plugin turns the app into an location tracking app for the time it runs in the background.