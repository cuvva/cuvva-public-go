package test_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cuvva/cuvva-public-go/lib/crpc"
)

// ExampleRequest_Headers demonstrates how to read and manipulate HTTP headers
// in CRPC request handlers.
func ExampleRequest_Headers() {
	// Create a new CRPC server
	server := crpc.NewServer(func(next crpc.HandlerFunc) crpc.HandlerFunc {
		return func(res http.ResponseWriter, req *crpc.Request) error {
			// Simple authentication middleware - just pass through
			return next(res, req)
		}
	})

	// Register a handler that demonstrates header usage
	server.Register("header_demo", "preview", nil, func(ctx context.Context) (*struct {
		Message string `json:"message"`
	}, error) {
		// Get the request from context
		req := crpc.GetRequestContext(ctx)
		if req == nil {
			return nil, fmt.Errorf("no request in context")
		}

		// Read headers from the request
		userAgent := req.GetHeader("User-Agent")
		authorization := req.GetHeader("Authorization")
		customHeader := req.GetHeader("X-Custom-Header")

		// You can also modify headers (though this is less common)
		req.SetHeader("X-Processed", "true")
		req.AddHeader("X-Debug", "header-demo")

		// Return response with header information
		return &struct {
			Message string `json:"message"`
		}{
			Message: fmt.Sprintf("User-Agent: %s, Auth: %s, Custom: %s", userAgent, authorization, customHeader),
		}, nil
	})

	// Create a test request with headers
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/preview/header_demo", nil)
	r.Header.Set("User-Agent", "example-client/1.0")
	r.Header.Set("Authorization", "Bearer token123")
	r.Header.Set("X-Custom-Header", "custom-value")

	// Handle the request
	server.ServeHTTP(w, r)

	fmt.Printf("Status: %d\n", w.Code)
	fmt.Printf("Response: %s", w.Body.String())

	// Output:
	// Status: 200
	// Response: {"message":"User-Agent: example-client/1.0, Auth: Bearer token123, Custom: custom-value"}
}
