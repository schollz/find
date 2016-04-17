# Fingerprint

Any computer with a WiFi card can use FIND using a client program that generates fingerprints by scanning the onboard wifi.

# Install

```bash
git clone https://github.com/schollz/find.git
cd find
git checkout fingerprint
go get ./...
go build
```

# Usage

```
./fingerprint
```

## Options

```
NAME:
   fingerprint - client for sending WiFi fingerprints to a FIND server

USAGE:
   find2.exe [global options] command [command options] [arguments...]

VERSION:
   0.2

COMMANDS:
GLOBAL OPTIONS:
   --server, -s "https://ml.internalpositioning.com"    server to connect
   --group, -g "group"                                  group name
   --user, -u "user"                                    user name
   --location, -l "location"                            location (needed for '--learn')
   --continue, -c "3"                                   number of times to run
   --learn, -e                                          need to set if you want to learn location
   --nodebug, -d                                        turns off debugging
   --help, -h                                           show help
   --version, -v                                        print the version
```

# Contributing

Currently seeking pull requests to add OS X and Windows support (see Issues).

# Credits

# License

Copyright 2015-2016 Zack Scholl

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License [https://github.com/schollz/find/blob/master/LICENSE](https://github.com/schollz/find/blob/master/LICENSE).

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
