# FIND API
[![](https://raw.githubusercontent.com/schollz/find/master/static/splash.gif)](https://www.internalpositioning.com/)

# Fingerprinting
## `POST` `/track` and `/learn`
Parses and inserts fingerprints. `time` is optional. `location` is optional for the `/track` route.

### POST

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
  "message": "Fingerprint inserted."
}
```

### curl

```bash
curl -H "Content-Type: application/json" -X POST -d 'JSON' http://server/learn
```

## `GET` `/calculate?group=X`
Recalculates the priors for the database for the `group`.

## `GET` `/location?group=X&user=Y&history=Z`
Gets the locations. If `user` is not provided it will return locations for all users in the `group`. If `history` is not included, it will return the last location, otherwise it will return the last `Z` locations.

### Response for `/location?group=something&user=user1&history=3`

```json
{
  "success": true,
  "message": "Successfully acquired.",
  "user1":{
      "time":"Some Date, 2010",
      "location":"location1",
      "bayes":{
        "location1":3.0,
        "location2":1.0,
        "location3":2.0,
        "location4":2.5,
      },
      "history":[
        {
          "location":"location1",
          "time":"time1"
        },
        {
          "location":"location1",
          "time":"time2"
        },
        {
          "location":"location2",
          "time":"time3"
        }
      ]
    }
}
```

# database

## `DELETE` `/username?group=X&user=Y`

Deletes user `Y` in group `X`.

### Response

```json
{
  "success":true,
  "message":"Deleted user Y"
}
```

# Meta

## `GET` `/status`

### Response


```json
{
  "num_cores":1,
  "registered":"2016-04-16 14:55:34.483803377 -0400 EDT",
  "status":"standard",
  "uptime":35109.225647597
}
```
