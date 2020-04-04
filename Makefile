all: plugins build

build:
	go build -v

rpi:
	GOOS=linux GOARCH=arm GOARM=7 go build -v

nano:
	export GOOS=linux GOARCH=arm GOARM=7 go build -v

plugins:
	make -C p
