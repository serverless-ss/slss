# options
ignore_output = &> /dev/null

# commands
install_apex = curl https://raw.githubusercontent.com/apex/apex/master/install.sh | sh

test:
	go test -v --race

apex:
	which apex $(ignore_output) || install_apex
