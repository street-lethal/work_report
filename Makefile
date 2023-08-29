help:
	@grep "^[a-zA-Z\-]*:" Makefile | grep -v "grep" | sed -e 's/^/make /' | sed -e 's/://'

test:
	goenv exec go test -v ./...

build:
	GOOS=linux GOARCH=amd64 goenv exec go build -trimpath -ldflags '-s -w' -o bin/main_linux
	GOOS=darwin GOARCH=arm64 goenv exec go build -trimpath -ldflags '-s -w' -o bin/main_mac

fmt:
	goenv exec go fmt ./...

lint:
	goenv exec go vet ./...

gen: ## Jira の HTML からレポート用データ(data/report.json)生成
	@./main -mode gen

send: ## プラットフォームに自動ログインしてレポート用データを送信
	@./main -mode send

send-m: ## 自動ログインせず、セッション情報(data/platform_session.json)を手動で設定してから送信
	@./main -mode send-m

clear: ## レポート用データを空データに更新する
	@./main -mode clear

gen-s:
	@goenv exec go run main.go -mode gen

send-s:
	@goenv exec go run main.go -mode send

send-m-s:
	@goenv exec go run main.go -mode send-m

clear-s:
	@goenv exec go run main.go -mode clear

init:
	cp config/settings.sample.json config/settings.json
	cp config/platform_id_test.json config/platform_id.json
	cp data/platform_session_test.json data/platform_session.json
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
