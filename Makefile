SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=findclient

VERSION=0.6
BUILD_TIME=`date +%FT%T%z`
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.VersionNum=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go get github.com/stretchr/testify/assert
	go get github.com/codegangsta/cli
	go get github.com/op/go-logging
	go build ${LDFLAGS} -o ${BINARY} ${SOURCES}

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm -rf builds
	rm -rf find
	rm -rf findclient*

.PHONY: binaries
binaries:
	go test
	rm -rf builds
	mkdir builds
	# Build Windows
	env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o findclient.exe -v *.go
	zip -r findclient_${VERSION}_windows_amd64.zip findclient.exe LICENSE
	mv findclient_${VERSION}_windows_amd64.zip builds/
	rm findclient.exe
	# Build Linux
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o findclient -v *.go
	zip -r findclient_${VERSION}_linux_amd64.zip findclient LICENSE
	mv findclient_${VERSION}_linux_amd64.zip builds/
	rm findclient
	# Build OS X
	env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o findclient -v *.go
	zip -r findclient_${VERSION}_osx.zip findclient LICENSE
	mv findclient_${VERSION}_osx.zip builds/
	rm findclient
	# Build Raspberry Pi / Chromebook
	env GOOS=linux GOARCH=arm go build ${LDFLAGS} -o findclient -v *.go
	zip -r findclient_${VERSION}_linux_arm.zip findclient LICENSE
	mv findclient_${VERSION}_linux_arm.zip builds/
	rm findclient
