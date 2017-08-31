# options
ignore_output = &> /dev/null

# commands
install_apex = curl https://raw.githubusercontent.com/apex/apex/master/install.sh | sh
get_shadowsocks = git clone --depth 1  https://github.com/shadowsocks/shadowsocks-go

test:
	go test -v --race

ensure_apex:
	@which apex $(ignore_output) || $(install_apex)

ensure_shadowsocks:
	@mkdir -p ./lambda/bin
	@ls shadowsocks-go $(ignore_output) || $(get_shadowsocks)
	GOOS=linux GOARCH=amd64 go build -o ./lambda/functions/slss/bin/shadowsocks_server ./shadowsocks-go/cmd/shadowsocks-server
	go build -o ./bin/shadowsocks_local ./shadowsocks-go/cmd/shadowsocks-local

install: ensure_apex ensure_shadowsocks
	go build -o ./bin/slss ./cmd/main.go

clean_up:
	@rm -rf ./bin
	@rm -rf ./shadowsocks-go
	@rm -rf ./lambda/functions/slss/bin

.PHONY: test ensure_apex ensure_shadowsocks install clean_up
