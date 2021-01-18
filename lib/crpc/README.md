# crpc

crpc contains a client and server conforming to the [Cuvva RPC standard](https://github.com/cuvva/standards).

It is heavily influenced by [net/rpc](https://golang.org/pkg/net/rpc/) and [Monzo Typhon](https://godoc.org/github.com/monzo/typhon).


## Components

crpc consists of a Client and a Server component.


### Client

The Client component is not intended to be used directly, but to be composed into a more fully-featured service client.

See [example/client/](/example/client) for example usage.


### Server

The Server component is intended to be used directly, and have handlers associated directly with it.

It implements `net/http.Handler`, thus can be embedded directly within an HTTP server. This is in preparation of enabling TLS between service, and thus internal RPC can use HTTP/2 multiplexing.

See [example/server/](/example/server/) for example usage.
