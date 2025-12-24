# - Host OS Platform
ifeq ($(OS),Windows_NT) 
    detected_OS := Windows
else
    detected_OS := $(sh -c 'uname 2>/dev/null || echo Unknown')
endif

# - Go Build Environment
GO=go
GO_TAGS=-tags mlog
ifeq ($(detected_OS), Windows)
	GOFLAGS = -v -buildmode=exe -gcflags all=-N 
	EXE_EXT=.exe
else
	GOFLAGS = -v -buildmode=pie
	EXE_EXT=
endif

# - Source Project Environment
# get the Makefile's directory (GNU Make >= v3.81)
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))
# set the GO Project's BIN directory
GO_PROJ_BIN=${mkfile_dir}bin

# - Packagers only
PKG_FULL_VERSION=$(shell grep -m 1 'MANUAL_VERSION' version.go | sed -E 's/.*"([^"]+)".*/\1/')
PKG_PUBLIC_NAME=$(shell grep -m 1 'module' go.mod | sed -E 's/^module\s+//p')
PKG_NAME=caesardisk
PKG_REVISION=1
PKG_VERSION=1.2
PKG_ARCH=amd64
PKG_FULLNAME=${PKG_NAME}_${PKG_VERSION}-${PKG_REVISION}_${PKG_ARCH}
PKG_BUILD_DIR=${HOME}/Develop/Distrib/Build/${PKG_NAME}
PKG_PPA_DIR=${HOME}/Develop/Distrib/PPA

# - Application stanza
EXEC_CAESAR=caesardisk
MAIN_CAESAR=cmd/disk/*.go
BIN_OUT_1=$(GO_PROJ_BIN)/$(EXEC_CAESAR)$(EXE_EXT)

# - Main Targets
.PHONY: clean build

all: caesardisk
	
allwin:
	GOOS=windows GOARCH=amd64 GOWORK=off $(GO) build $(GO_TAGS) $(GOFLAGS) -o ${BIN_OUT_1}.exe ${MAIN_CAESAR}

release:
	strip --strip-unneeded ${BIN_OUT_1}

version:
	@echo $(PKG_FULL_VERSION)

proxy:
	GOPROXY=proxy.golang.org go list -m $(PKG_PUBLIC_NAME)@v$(PKG_FULL_VERSION)

# - Application Targets

caesardisk:
	$(GO) build $(GO_TAGS) $(GOFLAGS) -o ${BIN_OUT_1} ${MAIN_CAESAR}

# - Secondary Targets

clean:
	go clean

run:
	go run -race  $(MAIN)

lint: 
	@gofmt -l . | grep ".*\.go"

test:
	go test tests/*test.go	

testall:
	go test ./...	

update:
	go get -u all

help:
	@echo "· Application related"
	@echo  "\tall - make ALL application targets on/for Linux"
	@echo  "\tallwin - make ALL application targets on/for Windows"
	@echo  "\tversion - print the application version found in the source code"
	@echo "· GO Language targets"
	@echo  "\ttest - Run all the tests"
	@echo  "\tupdate - Update all 3rd party GO package dependencies"
	@echo  "\tlint - Run GO Lint"
	@echo "· Package Building (NOT IMPLEMENTED)"
	@echo  "\tdebian - Build the Debian (${PKG_FULLNAME}-${PKG_REVISION}.deb) package"
	@echo  "\trpm - Build the RPM (${PKG_FULLNAME}.rpm) package"
	@echo  "\trpmclean - Cleans the RPM build area"

# Package Builders
