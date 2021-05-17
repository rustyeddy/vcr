all: plugins redeye

redeye: *.go
	go build -v

test:
	./bin/get-health.sh

rpi:
	GOOS=linux GOARCH=arm GOARM=7 go build -v

nano:
	export GOOS=linux GOARCH=arm GOARM=7 go build -v

plugins:
	make -C plugins

.PHONY: all build rpi nano plugins
