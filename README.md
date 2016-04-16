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

```bash
NAME:
   fingerprint - client for sending WiFi fingerprints to a FIND server

USAGE:
   compfind [global options] command [command options] [arguments...]

VERSION:
   0.1

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --server, -s "https://ml.internalpositioning.com"    server to connect
   --group, -g "group"                                  group name
   --user, -u "user"                                    user name
   --location, -l "location"                            location (needed for '--learn')
   --continue, -c "3"                                   number of times to run
   --learn, -e                                          need to set if you want to learn location
   --help, -h                                           show help
   --version, -v                                        print the version
```

# Contributing

Currently seeking pull requests to add OS X and Windows support (see Issues).

# Credits

# License

The MIT License (MIT).
