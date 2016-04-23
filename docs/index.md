# FIND API
[![](https://raw.githubusercontent.com/schollz/find/master/static/splash.gif)](https://www.internalpositioning.com/)

# Fingerprinting
## `POST /learn`

Parses and inserts fingerprints. `time` is optional. `location` is optional for the `/track` route.

Requires posting a `WifiFingerprint` JSON:

```json
{
   "group":"some group",
   "username":"some user",
   "location":"some place",
   "time":12309123,
   "wififingerprint":[
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

## `POST /track`

Parses and inserts fingerprints. `time` is optional. `location` is optional for the `/track` route.

Requires posting a `WifiFingerprint` JSON:

```json
{
   "group":"some group",
   "username":"some user",
   "location":"some place",
   "time":12309123,
   "wififingerprint":[
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


## `GET /calculate?group=X`
Recalculates the priors for the database for the `group`.


### Response

```json
{
  "message":"Parameters optimized",
  "success":true
}
```

## `GET /location?group=X&user=Y or /location?group=X or /location?group=X&user=Y,Z,W`

Gets the locations. If `user` is not provided it will return locations for all users in the `group`. If `history` is not included, it will return the last location, otherwise it will return the last `Z` locations.

### Response

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


## `DELETE` `/username?group=X&user=Y`

Deletes user `Y` in group `X`.

### Response

```json
{
  "success":true,
  "message":"Deleted user Y"
}
```

## `GET` `/status`

Returns status of the server and some information about the computer.

### Response


```json
{
  "num_cores":1,
  "registered":"2016-04-16 14:55:34.483803377 -0400 EDT",
  "status":"standard",
  "uptime":35109.225647597
}
```
