# API

**The server for FIND allows manipulation of fingerprints directly through these API routes. Most useful is likely the [/location](/#get-location) route which gathers the most recent location for a user.**

[![](https://raw.githubusercontent.com/schollz/find/master/static/splash.gif)](https://www.internalpositioning.com/)

<br><br><br><br><br>

## POST /learn

### Description

Submit a fingerprint to be used for learning the classification of the location. The information for the fingerprint is gathered from the WiFi client - either the App or the program.

### Parameters

#### POST

```json
{
   "group":"some group",
   "username":"some user",
   "location":"some place",
   "time":12309123,
   "wifi-fingerprint":[
      {
         "mac":"AA:AA:AA:AA:AA:AA",
         "rssi":-45
      },
      {
         "mac":"BB:BB:BB:BB:BB:BB",
         "rssi":-55
      }
   ]
}
```

### Response

```json
{
    "success": true,
    "message": "Inserted fingerprint containing 23 APs for zack at zakhome floor 2 office"
}
```

<br><br><br><br><br>

## POST /track

Submit a fingerprint to be used for classifying the location. The information for the fingerprint is gathered from the WiFi client - either the App or the program.

### Parameters

#### POST

```json
{
   "group":"some group",
   "username":"some user",
   "location":"some place",
   "time":12309123,
   "wifi-fingerprint":[
      {
         "mac":"AA:AA:AA:AA:AA:AA",
         "rssi":-45
      },
      {
         "mac":"BB:BB:BB:BB:BB:BB",
         "rssi":-55
     }
   ]
}
```

### Response

```json
{
    "success": true,
    "message": "Calculated location: zakhome floor 2 office",
    "location": "zakhome floor 2 office",
    "bayes": {
        "zakhome floor 1 kitchen": 0.07353831034486494,
        "zakhome floor 2 bedroom": -0.9283974092154644,
        "zakhome floor 2 office": 0.8548590988705993
    }
}
```

<br><br><br><br><br>

## GET /calculate

### Description

Recalculates the priors for the database for the `group`.

### Parameters

Name  | Location | Description                 | Required
----- | -------- | --------------------------- | --------
group | query    | Defines the unique group ID | yes

### Response

```json
{
  "message":"Parameters optimized",
  "success":true
}
```

<br><br><br><br><br>

## GET /location

### Description

Gets the locations for the specified user(s) in the specified group.

### Parameters

Name  | Location | Description                                             | Required
----- | -------- | ------------------------------------------------------- | --------
group | query    | Defines the unique group ID                             | yes
user  | query    | Specifies a user to get location                        | no
users | query    | Specifies multiple users `users=X,Y,Z` to get histories | no

### Response

If `user` or `users` are not specified, then the location of all users are returned.

```json
{
   "message":"Correctly found locations.",
   "success":true,
   "users":{
      "morpheus":[
         {
            "time":"2016-04-18 15:59:38.146929368 -0400 EDT",
            "location":"office",
            "bayes":{
               "bed bath":-1.0796283868148098,
               "bedroom":-0.3253323338565688,
               "car":-0.11084494825121938,
               "dining":-0.21592336935362944,
               "kitchen":0.7779455402822841,
               "living":-0.5328733505357962,
               "office":1.486656848529739
            }
         }
      ],
      "zack":[
         {
            "time":"2016-04-20 07:27:47.960140659 -0400 EDT",
            "location":"office",
            "bayes":{
               "bed bath":-1.028454724759723,
               "bedroom":0.1239023145100694,
               "car":-0.1493711750580678,
               "dining":-0.4237049232002753,
               "kitchen":0.6637176338607336,
               "living":-0.701636080467658,
               "office":1.515546955114921
            }
         }
      ]
   }
}
```

<br><br><br><br><br>

## DELETE /username

### Description

Deletes all the tracking fingerprints for specified user in the specified group.

### Parameters

Name  | Location | Description                 | Required
----- | -------- | --------------------------- | --------
group | query    | Defines the unique group ID | yes
user  | query    | Specifies a user            | yes

### Response

```json
{
  "success":true,
  "message":"Deleted user Y"
}
```

<br><br><br><br><br>

## DELETE /locations

### Description

Bulk delete locations

### Parameters

Name  | Location | Description                                               | Required
----- | -------- | --------------------------------------------------------- | --------
group | query    | Defines the unique group ID                               | yes
names | query    | Enter locations seperated by commas, e.g. locations=X,Y,Z | yes

### Response

```json
{
  "success":true,
  "message":"Deleted X locations"
}
```

<br><br><br><br><br>

## DELETE /database

### Description

Delete database and all associated data.

### Parameters

Name  | Location | Description                 | Required
----- | -------- | --------------------------- | --------
group | query    | Defines the unique group ID | yes

### Response

```json
{
  "success":true,
  "message":"Successfully deleted X."
}
```

<br><br><br><br><br>

## PUT /mixin

### Description

Allows overriding of the `Mixin` parameter. Value of `0` uses only the RSSI Priors, while value of `1` uses only the Mac prevalence statistics.

### Parameters

Name  | Location | Description                                                       | Required
----- | -------- | ----------------------------------------------------------------- | --------
group | query    | Defines the unique group ID                                       | yes
mixin | query    | Specifiy a value between 0 and 1 to activate, or -1 to deactivate | yes

### Response

```json
{
  "message":"Overriding mixin for testdb, now set to 1",
  "success":true
}
```

<br><br><br><br><br>

## PUT /database

### Description

Migrate a database. This copies all the contents of one database to another. If the group does not exist, it will be created. The group that is migrated from is not deleted.

### Parameters

Name | Location | Description                                       | Required
---- | -------- | ------------------------------------------------- | --------
from | query    | Defines the unique group to migrate from          | yes
to   | query    | Defines the unique group to migrate database into | yes

### Response

```json
{
  "message":"Successfully migrated X to Y",
  "success":true
}
```

<br><br><br><br><br>

## PUT /mqtt

### Description

Allows you to access MQTT streams of your data. This is available on the public server using the 3rd party [mosquitto server](https://mosquitto.org/), if you want to setup, [see this documentation](https://doc.internalpositioning.com/mqtt/).

### Parameters

Name  | Location | Description                 | Required
----- | -------- | --------------------------- | --------
group | query    | Defines the unique group ID | yes

### Response

```json
{
    "message": "You have successfully set your password.",
    "password": "YOURPASSWORD",
    "success": true
}
```

<br><br><br><br><br>

## GET /status

### Description

Returns status of the server and some information about the computer.

### Parameters

None.

### Response

```json
{
  "num_cores":1,
  "registered":"2016-04-16 14:55:34.483803377 -0400 EDT",
  "status":"standard",
  "uptime":35109.225647597
}
```
