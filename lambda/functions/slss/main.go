package main

import (
	"encoding/json"
	"fmt"

	"github.com/apex/go-apex"
	"github.com/serverless-ss/slss"
)

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		config := new(slss.LambdaShadowSocksConfig)
		if err := json.Unmarshal(event, config); err != nil {
			return nil, err
		}

		fmt.Printf("config: %+v\n", config)

		return nil, nil
	})
}
