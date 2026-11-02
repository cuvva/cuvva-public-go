# CRPC Header Support

The CRPC `Request` struct now supports reading and manipulating HTTP headers from incoming requests.

## Reading Headers

You can read headers from the request using the `GetHeader` method:

```go
func myHandler(ctx context.Context) (*MyResponse, error) {
    req := crpc.GetRequestContext(ctx)
    
    // Read specific headers
    userAgent := req.GetHeader("User-Agent")
    authorization := req.GetHeader("Authorization")
    customHeader := req.GetHeader("X-Custom-Header")
    
    // Headers are case-insensitive
    contentType := req.GetHeader("content-type") // Same as "Content-Type"
    
    // Returns empty string if header doesn't exist
    missing := req.GetHeader("Non-Existent-Header") // Returns ""
    
    // ... rest of handler logic
}
```

## Setting Headers

You can set headers on the request using the `SetHeader` method:

```go
func myHandler(ctx context.Context) (*MyResponse, error) {
    req := crpc.GetRequestContext(ctx)
    
    // Set a header (replaces any existing value)
    req.SetHeader("X-Processed", "true")
    req.SetHeader("X-Handler", "myHandler")
    
    // ... rest of handler logic
}
```

## Adding Headers

You can add multiple values to a header using the `AddHeader` method:

```go
func myHandler(ctx context.Context) (*MyResponse, error) {
    req := crpc.GetRequestContext(ctx)
    
    // Add values to a header (appends to existing values)
    req.AddHeader("X-Debug", "step1")
    req.AddHeader("X-Debug", "step2")
    // Now X-Debug header has both "step1" and "step2" values
    
    // ... rest of handler logic
}
```

## Deleting Headers

You can remove headers using the `DelHeader` method:

```go
func myHandler(ctx context.Context) (*MyResponse, error) {
    req := crpc.GetRequestContext(ctx)
    
    // Remove a header completely
    req.DelHeader("X-Sensitive-Data")
    
    // ... rest of handler logic
}
```

## Direct Header Access

For advanced use cases, you can access the underlying `http.Header` directly:

```go
func myHandler(ctx context.Context) (*MyResponse, error) {
    req := crpc.GetRequestContext(ctx)
    
    // Access the raw http.Header
    if req.Header != nil {
        // Get all values for a header
        debugValues := req.Header["X-Debug"]
        
        // Check if header exists
        if _, exists := req.Header["Authorization"]; exists {
            // Header exists
        }
        
        // Iterate over all headers
        for name, values := range req.Header {
            for _, value := range values {
                // Process each header value
            }
        }
    }
    
    // ... rest of handler logic
}
```

## Important Notes

1. **Case Insensitivity**: Header names are case-insensitive. The methods use Go's standard `textproto.CanonicalMIMEHeaderKey` for canonicalization.

2. **Nil Safety**: All header methods are safe to call even when the `Header` field is nil. They will initialize the header map as needed.

3. **Request Context**: Headers are available through the request context using `crpc.GetRequestContext(ctx)`.

4. **Backward Compatibility**: The existing `BrowserOrigin` field continues to work as before and contains the value of the "Origin" header.

## Example

See `example_headers_test.go` for a complete working example of header usage in CRPC handlers.
