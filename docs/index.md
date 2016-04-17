# FIND API


<a href="https://www.internalpositioning.com/"><img src="https://raw.githubusercontent.com/schollz/find/master/static/splash.gif"></a>

# Fingerprinting

## `POST` `/track` and `/learn`

`time` is optional.
`location` is optional for the `/track` route.

### JSON

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

### curl

`curl -H "Content-Type: application/json" -X POST -d 'JSON' http://server/learn`
