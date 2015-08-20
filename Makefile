run:
	gofmt -w *.go
	go run *.go --bind=:9999 --bindinfo=:9991
build:
	go build -o $(GOPATH)/bin/osx-echo *.go;
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $(GOPATH)/bin/echo *.go;

