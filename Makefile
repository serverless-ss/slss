# options
ignore_output = &> /dev/null

# commands
install_apex = curl https://raw.githubusercontent.com/apex/apex/master/install.sh | sh

test:
	go test -v --race

ensure_apex:
	@which apex $(ignore_output) || $(install_apex)

.PHONY: test ensure_apex
