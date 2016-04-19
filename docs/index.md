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
  "message": "Inserted X fingerprints for USER at LOCATION."
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
  "message": "Calculated location: LOCATION",
  "location": "LOCATION"
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

## `GET /location?group=X&user=Y&history=Z`
Gets the locations. If `user` is not provided it will return locations for all users in the `group`. If `history` is not included, it will return the last location, otherwise it will return the last `Z` locations.

### Response

```json
{
    "success":true,
    "message":"Found X users.",
    "userX": [
        {
            "time": "2016-04-16 08:16:43.123233725 -0400 EDT",
            "location": "some location",
            "bayes": {}
        }
    ]
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
