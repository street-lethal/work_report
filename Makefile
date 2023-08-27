test:
	goenv exec go test -v ./...

build:
	GOOS=linux GOARCH=amd64 goenv exec go build -trimpath -ldflags '-s -w' -o bin/main_linux
	GOOS=darwin GOARCH=arm64 goenv exec go build -trimpath -ldflags '-s -w' -o bin/main_mac

fmt:
	goenv exec go fmt ./...

lint:
	goenv exec go vet ./...

gen:
	@./main -mode gen

send:
	@./main -mode send

clear:
	@./main -mode clear

gen-s:
	@goenv exec go run main.go -mode gen

send-s:
	@goenv exec go run main.go -mode send

clear-s:
	@goenv exec go run main.go -mode clear

init:
	cp config/settings.sample.json config/settings.json
	touch data/jira.html

init-linux:
	if [ -e main ]; then \
		unlink main; \
	fi
	ln -s bin/main_linux main

init-mac:
	if [ -e main ]; then \
		unlink main; \
	fi
	ln -s bin/main_mac main
