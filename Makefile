SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=findserver

VERSION=2.2
BUILD_TIME=`date +%FT%T%z`
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.VersionNum=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go get github.com/boltdb/bolt
	go get github.com/gin-gonic/contrib/sessions
	go get github.com/gin-gonic/gin
	go get github.com/stretchr/testify/assert
	go get github.com/pquerna/ffjson/fflib/v1
	go get github.com/eclipse/paho.mqtt.golang.git
	go build ${LDFLAGS} -o ${BINARY} ${SOURCES}

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm -rf builds
	rm -rf find
	rm -rf findserver*

.PHONY: binaries
binaries:
	rm -rf builds
	mkdir builds
	# Build Windows
	env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o findserver.exe -v *.go
	zip -r find_${VERSION}_windows_amd64.zip findserver.exe LICENSE ./templates/* ./data/.datagoeshere ./static/*
	mv find_${VERSION}_windows_amd64.zip builds/
	rm findserver.exe
	# Build Linux
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o findserver -v *.go
	zip -r find_${VERSION}_linux_amd64.zip findserver LICENSE ./templates/* ./data/.datagoeshere ./static/*
	mv find_${VERSION}_linux_amd64.zip builds/
	rm findserver
	# Build OS X
	env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o findserver -v *.go
	zip -r find_${VERSION}_osx.zip findserver LICENSE ./templates/* ./data/.datagoeshere ./static/*
	mv find_${VERSION}_osx.zip builds/
	rm findserver
	# Build Raspberry Pi / Chromebook
	env GOOS=linux GOARCH=arm go build ${LDFLAGS} -o findserver -v *.go
	zip -r find_${VERSION}_linux_arm.zip findserver LICENSE ./templates/* ./data/.datagoeshere ./static/*
	mv find_${VERSION}_linux_arm.zip builds/
	rm findserver
