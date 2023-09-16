.PHONY: install build run

APP_NAME := goTyper.exe

intall:
	go get

build:
	go build -o bin/$(APP_NAME) 

run: build
	./bin/$(APP_NAME)