NAME=plonker
VERSION=0.1.0

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

darwin_%: GOOS=darwin

linux_%: GOOS=linux

windows_%: GOOS=windows
windows_%: EXT=.exe

all: clean darwin_package linux_package windows_package

clean:
	rm -rf ./build

%_build:
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-s -w" -o ./build/$(GOOS)/$(GOARCH)/$(NAME)$(EXT)

%_package: %_build
	@cd ./build/$(GOOS)/$(GOARCH) && tar -czf ./$(NAME)-$(GOOS)-$(GOARCH)-$(VERSION).tar.gz ./$(NAME)$(EXT)
	@shasum -a256 ./build/$(GOOS)/$(GOARCH)/$(NAME)-$(GOOS)-$(GOARCH)-$(VERSION).tar.gz
