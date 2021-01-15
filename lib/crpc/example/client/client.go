package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/crpc"
	"github.com/cuvva/cuvva-public-go/lib/crpc/example"
	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
)

type ExampleClient struct {
	*crpc.Client
}

func (ec *ExampleClient) Ping(ctx context.Context) error {
	return ec.Do(ctx, "ping", "2017-11-08", nil, nil)
}

func (ec *ExampleClient) Greet(ctx context.Context, req *example.GreetRequest) (res *example.GreetResponse, err error) {
	return res, ec.Do(ctx, "greet", "2017-11-08", req, &res)
}

func main() {
	client := &http.Client{
		Transport: jsonclient.NewAuthenticatedRoundTripper(nil, "Bearer", "...someJWTOrSomething"),
		Timeout:   5 * time.Second,
	}

	var ec example.Service = &ExampleClient{
		Client: crpc.NewClient("http://127.0.0.1:3000/v1", client),
	}

	ctx := context.Background()

	if err := ec.Ping(ctx); err != nil {
		fmt.Printf("ping failed: %#v\n", err)
		return
	}

	res, err := ec.Greet(ctx, &example.GreetRequest{Name: "James"})
	if err != nil {
		fmt.Printf("could not greet: %#v\n", err)
		return
	}

	fmt.Println("greeting:", res.Greeting)
}
