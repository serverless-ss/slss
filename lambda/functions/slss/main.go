package main

import (
	"encoding/json"

	"github.com/apex/go-apex"
)

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		return nil, nil
	})
}
