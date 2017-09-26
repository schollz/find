SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=fingerprint

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
	rm -rf fingerprint*

.PHONY: binaries
binaries:
	rm -rf builds
	mkdir builds
	# Build Windows
	env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o fingerprint.exe -v *.go
	zip -r fingerprint_${VERSION}_windows_amd64.zip fingerprint.exe LICENSE
	mv fingerprint_${VERSION}_windows_amd64.zip builds/
	rm fingerprint.exe
	# Build Linux
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o fingerprint -v *.go
	zip -r fingerprint_${VERSION}_linux_amd64.zip fingerprint LICENSE
	mv fingerprint_${VERSION}_linux_amd64.zip builds/
	rm fingerprint
	# Build OS X
	env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o fingerprint -v *.go
	zip -r fingerprint_${VERSION}_osx.zip fingerprint LICENSE
	mv fingerprint_${VERSION}_osx.zip builds/
	rm fingerprint
	# Build Raspberry Pi / Chromebook
	env GOOS=linux GOARCH=arm go build ${LDFLAGS} -o fingerprint -v *.go
	zip -r fingerprint_${VERSION}_linux_arm.zip fingerprint LICENSE
	mv fingerprint_${VERSION}_linux_arm.zip builds/
	rm fingerprint
