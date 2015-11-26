# generate version number
version=$(shell git describe --tags --long --always|sed 's/^v//')

all: dep
	gom exec go build  -ldflags "-X main.version=$(version)" powerdns-remote.go
	go fmt


dep:
	gom install

version:
	@echo $(version)
