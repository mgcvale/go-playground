
build:
	go build -o bin/awesomeProject ./cmd/server.go

run:
	./bin/awesomeProject

dev:
	go build -o bin/awesomeProject ./cmd/server.go; ./bin/awesomeProject


install:
	go build -o /usr/bin/awesomeProject cmd/server.go

uninstall:
	rm /usr/bin/awesomeProject
