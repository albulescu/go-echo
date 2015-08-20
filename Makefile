install:
	go get gopkg.in/mgo.v2;
	go get gopkg.in/mgo.v2/bson;
	go get gopkg.in/ini.v1
	go get github.com/dvsekhvalnov/jose2go
run:
	gofmt -w *.go
	go run *.go --config=config.ini
build:
	gofmt -w *.go
	go build -o $(GOPATH)/bin/osx-echo *.go;
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $(GOPATH)/bin/echo *.go;

