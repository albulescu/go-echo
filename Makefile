run:
	go run *.go
build:
	go build -o $(GOPATH)/bin/osx-echo *.go;
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $(GOPATH)/bin/echo *.go;
