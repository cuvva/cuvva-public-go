package jsonclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	pathlib "path"
	"strings"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/middleware/request"
	"github.com/cuvva/cuvva-public-go/lib/version"
)

var (
	// ErrNoResponse is returned when a client request is given a body to
	// unmarshal to however the server does not return any content (HTTP 204).
	ErrNoResponse = &ClientRequestError{"no response to unmarshal to body", nil}
)

type KeyVersion string

// Version1 is an auth key version that is provided via config file
const Version1 KeyVersion = "01."

// DefaultUserAgent is the default HTTP User-Agent Header that is presented to the server.
var DefaultUserAgent = "jsonclient/" + version.Truncated + " (+https://cuvva.com)"

// Client represents a json-client HTTP client.
type Client struct {
	Scheme string
	Host   string
	Prefix string

	UserAgent string

	Client *http.Client
}

// NewClient returns a client configured with a transport scheme, remote host
// and URL prefix supplied as a URL <scheme>://<host></prefix>
func NewClient(baseURL string, c *http.Client) *Client {
	remote, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	if c == nil {
		c = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return &Client{
		Scheme: remote.Scheme,
		Host:   remote.Host,
		Prefix: remote.Path,

		UserAgent: DefaultUserAgent,

		Client: c,
	}
}

// Do executes an HTTP request against the configured server.
func (c *Client) Do(ctx context.Context, method, path string, params url.Values, src, dst interface{}) error {
	if c.Client == nil {
		c.Client = http.DefaultClient
	}

	if ctx == nil {
		ctx = context.Background()
	}

	ctx, requestID := request.GetOrSetRequestID(ctx)

	req := &http.Request{
		Method: method,
		URL: &url.URL{
			Scheme: c.Scheme,
			Host:   c.Host,
			Path:   pathlib.Join(c.Prefix, path),
		},
		Header: http.Header{
			"Accept":     []string{"application/json"},
			"User-Agent": []string{c.UserAgent},
			"Request-Id": []string{requestID},
		},
		Host: c.Host,
	}

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	err := c.setRequestBody(req, src)
	if err != nil {
		return &ClientRequestError{"could not marshal", err}
	}

	res, err := c.Client.Do(req.WithContext(ctx))
	if err != nil {
		if netErr, ok := err.(net.Error); ok {
			if netErr.Timeout() {
				return cher.New(cher.RequestTimeout, cher.M{"method": method, "path": path, "host": c.Host, "scheme": c.Scheme})
			}

			return &ClientTransportError{method, path, "request failed", netErr}
		}

		return &ClientTransportError{method, path, "unknown error", err}
	}

	defer res.Body.Close()

	return c.handleResponse(res, method, path, dst)
}

func (c *Client) setRequestBody(req *http.Request, src interface{}) error {
	if src != nil {
		var buf bytes.Buffer

		err := json.NewEncoder(&buf).Encode(src)
		if err != nil {
			return err
		}

		req.Body = ioutil.NopCloser(&buf)
		req.ContentLength = int64(buf.Len())

		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	return nil
}

func (c *Client) handleResponse(res *http.Response, method, path string, dst interface{}) error {
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		if dst == nil {
			return nil
		}

		if res.StatusCode == http.StatusNoContent || res.Body == nil {
			return ErrNoResponse
		}

		err := json.NewDecoder(res.Body).Decode(dst)
		if err == io.EOF {
			return ErrNoResponse
		} else if err != nil {
			return &ClientTransportError{method, path, "could not unmarshal", err}
		}

		return nil
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &ClientTransportError{method, path, "could not read response body stream", err}
	}

	var body cher.E
	if err := json.Unmarshal(resBody, &body); err == nil && body.Code != "" {
		return body
	}

	var errorResBody interface{}
	if err := json.Unmarshal(resBody, &errorResBody); err != nil {
		errorResBody = string(resBody)
	}

	statusText := http.StatusText(res.StatusCode)
	if statusText == "" {
		statusText = "unknown"
	}

	statusParts := strings.Fields(statusText)

	for i := range statusParts {
		statusParts[i] = strings.ToLower(statusParts[i])
	}

	newErrorMessage := strings.Join(statusParts, "_")

	return cher.New(newErrorMessage, cher.M{
		"httpStatus": res.StatusCode,
		"data":       errorResBody,
		"method":     res.Request.Method,
		"url":        res.Request.URL.String(),
	})
}

// ClientRequestError is returned when an error related to
// constructing a client request occurs.
type ClientRequestError struct {
	ErrorString string

	cause error
}

// Cause returns the causal error (if wrapped) or nil
func (cre *ClientRequestError) Cause() error {
	return cre.cause
}

func (cre *ClientRequestError) Error() string {
	if cre.cause != nil {
		return cre.ErrorString + ": " + cre.cause.Error()
	}

	return cre.ErrorString
}

// ClientTransportError is returned when an error related to
// executing a client request occurs.
type ClientTransportError struct {
	Method, Path, ErrorString string

	cause error
}

// Cause returns the causal error (if wrapped) or nil
func (cte *ClientTransportError) Cause() error {
	return cte.cause
}

func (cte *ClientTransportError) Error() string {
	if cte.cause != nil {
		return fmt.Sprintf("%s %s %s: %s", cte.Method, cte.Path, cte.ErrorString, cte.cause.Error())
	}

	return fmt.Sprintf("%s %s %s", cte.Method, cte.Path, cte.ErrorString)
}

// AuthenticatedRoundTripper applies authentication before handing the request
// to the embedded transport for execution.
type AuthenticatedRoundTripper struct {
	http.RoundTripper

	authHeader string
}

// NewAuthenticatedRoundTripper returns a new AuthenticatedRoundTripper that will
// apply the given auth before performing the request.
func NewAuthenticatedRoundTripper(rt http.RoundTripper, authType, authToken string) *AuthenticatedRoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}

	if authToken == string(Version1) {
		panic("no authentication token added")
	}

	return &AuthenticatedRoundTripper{
		RoundTripper: rt,

		authHeader: authType + " " + authToken,
	}
}

// RoundTrip applies authentication before performing the request.
func (art *AuthenticatedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", art.authHeader)

	return art.RoundTripper.RoundTrip(req)
}

// BasicAuthRoundTripper applies authentication before handing the request
// to the embedded transport for execution.
type BasicAuthRoundTripper struct {
	http.RoundTripper

	username string
	password string
}

// NewBasicAuthRoundTripper returns a new BasicAuthRoundTripper that will
// apply the given auth before performing the request.
func NewBasicAuthRoundTripper(rt http.RoundTripper, username, password string) *BasicAuthRoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return &BasicAuthRoundTripper{
		RoundTripper: rt,

		username: username,
		password: password,
	}
}

// RoundTrip applies authentication before performing the request.
func (bart *BasicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(bart.username, bart.password)

	return bart.RoundTripper.RoundTrip(req)
}
