build:
	go build -v

rpi:
	export GOOS=linux GOARCH=arm GOARM=7 go build -v

nano:
	export GOOS=linux GOARCH=arm GOARM=7 go build -v
