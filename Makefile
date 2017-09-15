# options
ignore_output = &> /dev/null

# commands
install_apex = curl https://raw.githubusercontent.com/apex/apex/master/install.sh | sh
get_shadowsocks = git clone --depth 1 https://github.com/shadowsocks/shadowsocks-go
get_ngrok_darwin = wget https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-darwin-amd64.zip && unzip ./ngrok-stable-darwin-amd64.zip && mv ngrok ./bin/ngrok && rm -rf ngrok-stable-darwin-amd64.zip
get_ngrok_lambda = wget https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-amd64.zip && unzip ./ngrok-stable-linux-amd64.zip && mv ngrok ./lambda/functions/slss/bin/ngrok && rm -rf ngrok-stable-linux-amd64.zip

test:
	go test -v --race

ensure_apex:
	@which apex $(ignore_output) || $(install_apex)

ensure_shadowsocks:
	@mkdir -p ./bin
	@mkdir -p ./lambda/functions/slss/bin
	@ls shadowsocks-go $(ignore_output) || $(get_shadowsocks)
	GOOS=linux GOARCH=amd64 go build -o ./lambda/functions/slss/bin/shadowsocks_server ./shadowsocks-go/cmd/shadowsocks-server
	go build -o ./bin/shadowsocks_local ./shadowsocks-go/cmd/shadowsocks-local
	@rm -rf ./shadowsocks-go

ensure_ngrok:
	@mkdir -p ./lambda/functions/slss/bin
	@ls ./bin/ngrok $(ignore_output) || $(get_ngrok_darwin)
	@ls ./lambda/functions/slss/bin/ngrok $(ignore_output) || $(get_ngrok_lambda)

ensure_all: ensure_apex ensure_shadowsocks ensure_ngrok

clean_up:
	@rm -rf ./bin
	@rm -rf ./shadowsocks-go
	@rm -rf ./lambda/functions/slss/bin

install: ensure_all
	@go build -o ./bin/slss ./cmd/main.go

start:
	@./bin/slss -c ./config.json

.PHONY: test ensure_apex ensure_shadowsocks ensure_ngrok ensure_all install start clean_up
