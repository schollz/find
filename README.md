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

Fingerprint is a CLI client for the Framework for Internal Navigation and Discovery.

Copyright (C) 2015-2016 Zack Scholl

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the [GNU Affero General Public License](LICENSE) for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see [GNU Affero General Public License here](https://www.gnu.org/licenses/agpl.html).
