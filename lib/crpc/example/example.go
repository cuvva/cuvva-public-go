package example

import (
	"context"

	"github.com/xeipuuv/gojsonschema"
)

type Service interface {
	Ping(context.Context) error
	Greet(context.Context, *GreetRequest) (*GreetResponse, error)
}

type GreetRequest struct {
	Name string `json:"name"`
}

var GreetRequestSchema = gojsonschema.NewStringLoader(`{
	"$schema": "http://json-schema.org/schema#",

	"type": "object",
	"required": [ "name" ],

	"properties": {
		"name": {
			"type": "string"
		}
	}
}`)

type GreetResponse struct {
	Greeting string `json:"greeting"`
}
