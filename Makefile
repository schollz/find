SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=find

VERSION=2.2
BUILD_TIME=`date +%FT%T%z`
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.VersionNum=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	# go get github.com/boltdb/bolt
	# go get github.com/gin-gonic/contrib/sessions
	# go get github.com/gin-gonic/gin
	# go get github.com/stretchr/testify/assert
	# go get github.com/pquerna/ffjson/fflib/v1
	# go get git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git
	go build ${LDFLAGS} -o ${BINARY} ${SOURCES}

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm -rf binaries

.PHONY: binaries
binaries:
	rm -rf binaries
	mkdir binaries
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o binaries/sdees
	zip -j binaries/sdees_linux_amd64.zip binaries/sdees
	rm binaries/sdees
	env GOOS=linux GOARCH=arm go build ${LDFLAGS} -o binaries/sdees
	zip -j binaries/sdees_linux_arm.zip binaries/sdees
	rm binaries/sdees
	env GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o binaries/sdees
	zip -j binaries/sdees_linux_arm64.zip binaries/sdees
	rm binaries/sdees
	wget ftp://ftp.vim.org/pub/vim/pc/vim74w32.zip
	unzip vim74w32.zip
	mv vim/vim74/vim.exe ./binaries/
	env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o binaries/sdees.exe
	zip -j binaries/sdees_windows_amd64.zip binaries/sdees.exe binaries/vim.exe
	rm -rf binaries/vim.exe
	rm -rf ./vim/
	rm binaries/sdees.exe
