# slss
[![Build Status](https://travis-ci.org/serverless-ss/slss.svg?branch=master)](https://travis-ci.org/serverless-ss/slss)

A port of shadowsocks running in Amazon Lambda.

## How to use

### Get the source code

```
go get -v github.com/serverless-ss/slss
```

### Edit the `config.json`

```
cd $GOPATH/src/github.com/serverless-ss/slss
vim config.json
```

### Install the dependencies and run slss

```
make install && make start
```
