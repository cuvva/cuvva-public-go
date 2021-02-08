package crpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"sort"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/xeipuuv/gojsonschema"
)

// ResponseWriter is the destination for RPC responses.
type ResponseWriter interface {
	io.Writer
}

// Request contains metadata about the RPC request.
type Request struct {
	Version string
	Method  string

	Body io.ReadCloser

	RemoteAddr    string
	BrowserOrigin string

	ctx context.Context
}

// Context returns the requests context from the transport.
//
// The returned context is always non-nil, it defaults to the
// background context.
func (r *Request) Context() context.Context {
	if r.ctx == nil {
		return context.Background()
	}

	return r.ctx
}

// WithContext sets the context of a Request
func (r *Request) WithContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

type contextKey string

const requestKey contextKey = "crpcrequest"

// GetRequestContext returns the Request from the context object
func GetRequestContext(ctx context.Context) *Request {
	if val, ok := ctx.Value(requestKey).(*Request); ok {
		return val
	}

	return nil
}

func setRequestContext(ctx context.Context, request *Request) context.Context {
	return context.WithValue(ctx, requestKey, request)
}

// HandlerFunc defines a handler for an RPC request. Request and response body
// data will be JSON. Request will immediately io.EOF if there is no request.
type HandlerFunc func(res http.ResponseWriter, req *Request) error

// MiddlewareFunc is a function that wraps HandlerFuncs.
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

// WrappedFunc contains the wrapped handler, and some additional information
// about the function that was determined during the reflection process
type WrappedFunc struct {
	Handler       HandlerFunc
	AcceptsInput  bool
	ReturnsResult bool
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()
var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()

// Wrap reflects a HandlerFunc from any function matching the
// following signatures:
//
// func(ctx context.Context, request *T) (response *T, err error)
// func(ctx context.Context, request *T) (err error)
// func(ctx context.Context) (response *T, err error)
// func(ctx context.Context) (err error)
func Wrap(fn interface{}) (*WrappedFunc, error) {
	// prevent re-reflection of type that is already a HandlerFunc
	if _, ok := fn.(HandlerFunc); ok {
		return nil, fmt.Errorf("fn doesn't need to be wrapped, use RegisterFunc")
	} else if _, ok := fn.(*WrappedFunc); ok {
		return nil, fmt.Errorf("fn is already wrapped, use RegisterFunc")
	}

	v := reflect.ValueOf(fn)
	t := v.Type()

	// check the basic type and the number of inputs/outputs
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("fn must be function, got %s", t.Kind())
	} else if t.NumIn() < 1 || t.NumIn() > 2 {
		return nil, fmt.Errorf("fn input must be (context.Context) or (context.Context, *T), got %d arguments", t.NumIn())
	} else if t.NumOut() < 1 || t.NumOut() > 2 {
		return nil, fmt.Errorf("fn output must be (error) or (*T, error), got %d arguments", t.NumOut())
	}

	if !t.In(0).Implements(contextType) {
		return nil, fmt.Errorf("fn first argument must implement context.Context, got %s", t.In(0))
	} else if !t.Out(t.NumOut() - 1).Implements(errorType) {
		return nil, fmt.Errorf("fn last argument must implement error, got %s", t.Out(t.NumOut()-1))
	}

	// resolve function parameter pointers to underlying type for use with
	// reflect.New (which will return pointers).
	var reqT, resT reflect.Type = nil, nil

	if t.NumIn() == 2 {
		if t.In(1).Kind() != reflect.Ptr {
			return nil, fmt.Errorf("fn last argument must be a pointer, got %s", t.In(1))
		}

		reqT = t.In(1).Elem()
		if reqT.Kind() != reflect.Struct {
			return nil, fmt.Errorf("fn last argument must be a struct, got %s", reqT.Kind())
		}
	}

	if t.NumOut() == 2 {
		var err error

		resT, err = wrapReturn(t.Out(0))
		if err != nil {
			return nil, err
		}
	}

	hn := func(w http.ResponseWriter, r *Request) error {
		ctx := reflect.ValueOf(r.Context())
		var inputs []reflect.Value

		if reqT == nil {
			if r.Body != nil {
				i, err := r.Body.Read(make([]byte, 1))
				if i != 0 || err != io.EOF {
					return cher.New(cher.BadRequest, nil, cher.New("unexpected_request_body", nil))
				}
			}

			inputs = []reflect.Value{ctx}
		} else {
			if r.Body == nil {
				return cher.New(cher.BadRequest, nil, cher.New("missing_request_body", nil))
			}

			req := reflect.New(reqT)
			err := json.NewDecoder(r.Body).Decode(req.Interface())
			if err == io.EOF {
				return cher.New(cher.BadRequest, nil, cher.New("missing_request_body", nil))
			} else if err != nil {
				return fmt.Errorf("crpc: json decoder error: %w", err)
			}

			inputs = []reflect.Value{ctx, req}
		}

		res := v.Call(inputs)

		if err := res[len(res)-1]; !err.IsNil() {
			return err.Interface().(error)
		}

		if len(res) == 1 {
			w.WriteHeader(http.StatusNoContent)
		} else if len(res) == 2 {
			enc := json.NewEncoder(w)
			enc.SetEscapeHTML(false)
			err := enc.Encode(res[0].Interface())
			if err != nil {
				return err
			}
		}

		return nil
	}

	return &WrappedFunc{
		Handler:       hn,
		AcceptsInput:  reqT != nil,
		ReturnsResult: resT != nil,
	}, nil
}

func wrapReturn(t reflect.Type) (reflect.Type, error) {
	switch t.Kind() {
	case reflect.Ptr:
		n := t.Elem()

		if n.Kind() != reflect.Struct {
			_, err := wrapReturn(n)
			if err != nil {
				return nil, err
			}
		}

		return n, nil

	case reflect.Slice:
		n := t.Elem()

		if n.Kind() != reflect.String {
			_, err := wrapReturn(n)
			if err != nil {
				return nil, err
			}
		}

		return t, nil

	default:
		return nil, fmt.Errorf("unsupported return type, expected *struct or slice; got %s", t)
	}
}

// MustWrap is the same as Wrap, however it panics when passed an
// invalid handler.
func MustWrap(fn interface{}) *WrappedFunc {
	wrapped, err := Wrap(fn)
	if err != nil {
		panic(err)
	}

	return wrapped
}

type handler struct {
	v  string
	fn HandlerFunc
}

// Server is an HTTP-compatible crpc handler.
type Server struct {
	// AuthenticationMiddleware applies authentication before any other
	// middleware or request is processed. Servers without middleware must
	// configure the UnsafeNoAuthentication middleware. If no
	// AuthenticationMiddleware is configured, the server will panic.
	AuthenticationMiddleware MiddlewareFunc

	// methods = version -> method -> HandlerFunc
	registeredVersionMethods map[string]map[string]*handler
	registeredPreviewMethods map[string]*handler

	resolvedMethods map[string]map[string]*handler

	mw []MiddlewareFunc
}

// NewServer returns a new RPC Server with an optional exception tracker.
func NewServer(auth MiddlewareFunc) *Server {
	return &Server{
		AuthenticationMiddleware: auth,
	}
}

// Use includes a piece of Middleware to wrap HandlerFuncs.
func (s *Server) Use(mw MiddlewareFunc) {
	s.mw = append(s.mw, mw)
}

const (
	// VersionPreview is used for experimental endpoints in development which
	// are coming but a version identifier has not been decided yet or may
	// be withdrawn at any point.
	VersionPreview = "preview"

	// VersionLatest is used by engineers only to call the latest version
	// of an endpoint in utilities like cURL and Paw.
	VersionLatest = "latest"
)

// expVersion matches valid method versions
var expVersion = regexp.MustCompile(`^(?:preview|20\d{2}-\d{2}-\d{2})$`)

// expMethod matched valid method names
var expMethod = regexp.MustCompile(`^[a-z][a-z\d]*(?:_[a-z\d]+)*$`)

func isValidMethod(method, version string) bool {
	return expMethod.MatchString(method) && expVersion.MatchString(version)
}

// Register reflects a HandlerFunc from fnR and associates it with a
// method name and version. If fnR does not meet the HandlerFunc standard
// defined above, or the presence of the schema doesn't match the presence
// of the input argument, Register will panic. This function is not thread safe
// and must be run in serial if called multiple times.
func (s *Server) Register(method, version string, schema gojsonschema.JSONLoader, fnR interface{}, mw ...MiddlewareFunc) {
	s.RegisterValidated(method, version, schema, nil, fnR, mw...)
}

func (s *Server) RegisterValidated(method, version string, reqSchema gojsonschema.JSONLoader, respSchema gojsonschema.JSONLoader, fnR interface{}, mw ...MiddlewareFunc) {
	if fnR == nil {
		s.RegisterFunc(method, version, reqSchema, nil, mw...)

		return
	}

	wrapped := MustWrap(fnR)
	hasSchema := reqSchema != nil

	if wrapped.AcceptsInput != hasSchema {
		if hasSchema {
			panic("schema validation configured, but handler doesn't accept input")
		} else {
			panic("no schema validation configured")
		}
	}

	s.RegisterValidatedFunc(method, version, reqSchema, respSchema, &wrapped.Handler, mw...)
}

// RegisterFunc associates a method name and version with a HandlerFunc,
// and optional middleware. This function is not thread safe and must be run in
// serial if called multiple times.
func (s *Server) RegisterFunc(method, version string, reqSchema gojsonschema.JSONLoader, fn *HandlerFunc, mw ...MiddlewareFunc) {
	s.RegisterValidatedFunc(method, version, reqSchema, nil, fn, mw...)
}

func (s *Server) RegisterValidatedFunc(method, version string, reqSchema gojsonschema.JSONLoader, respSchema gojsonschema.JSONLoader, fn *HandlerFunc, mw ...MiddlewareFunc) {
	if s.registeredVersionMethods == nil {
		s.registeredVersionMethods = make(map[string]map[string]*handler)
	}

	if s.registeredPreviewMethods == nil {
		s.registeredPreviewMethods = make(map[string]*handler)
	}

	if fn == nil && reqSchema != nil {
		panic("schema validation configured, but handler is nil")
	}

	if !isValidMethod(method, version) {
		panic("invalid method/version")
	} else if s.AuthenticationMiddleware == nil {
		panic("no authentication configured")
	}

	if fn == nil && version == VersionPreview {
		panic(fmt.Sprintf("cannot set preview method '%s' as nil", method))
	}

	if s.isRouteDefined(method, version) {
		panic(fmt.Sprintf("cannot set '%s' on version '%s', it's already defined", method, version))
	}

	if fn == nil {
		s.setRoute(version, method, nil)
		s.buildRoutes()
		return
	}

	middleware := []MiddlewareFunc{s.AuthenticationMiddleware}
	if reqSchema != nil {
		compiledSchema, err := gojsonschema.NewSchemaLoader().Compile(reqSchema)
		if err != nil {
			panic(fmt.Sprintf("request schema error in %s: %s", method, err))
		}

		middleware = append(middleware, Validate(compiledSchema))
	}

	if respSchema != nil {
		compiledSchema, err := gojsonschema.NewSchemaLoader().Compile(respSchema)
		if err != nil {
			panic(fmt.Sprintf("response schema error in %s: %s", method, err))
		}

		middleware = append(middleware, ValidateResponseMiddleware(compiledSchema))
	}

	middleware = append(middleware, mw...)

	// This wraps the middleware funcs inside each one in reverse order
	for i := range middleware {
		p := mw[len(mw)-1-i](*fn)
		fn = &p
	}

	for i := range s.mw {
		p := s.mw[len(s.mw)-1-i](*fn)
		fn = &p
	}

	s.setRoute(version, method, &handler{version, *fn})
	s.buildRoutes()
}

func (s Server) isRouteDefined(method, version string) bool {
	if version == VersionPreview {
		_, ok := s.registeredPreviewMethods[method]
		return ok
	}

	if methodSet, ok := s.registeredVersionMethods[version]; ok {
		_, ok := methodSet[method]
		return ok
	}

	return false
}

func (s *Server) setRoute(version, method string, hn *handler) {
	if version == VersionPreview {
		s.registeredPreviewMethods[method] = hn

		return
	}

	versions, ok := s.registeredVersionMethods[version]
	if !ok {
		versions = make(map[string]*handler)
		s.registeredVersionMethods[version] = versions
	}

	versions[method] = hn
}

func (s *Server) buildRoutes() {
	knownVersions := sort.StringSlice{}
	resolvedMethods := make(map[string]map[string]*handler)

	// build known versions
	for version, methodSet := range s.registeredVersionMethods {
		if methodSet == nil {
			continue
		}

		if _, ok := resolvedMethods[version]; !ok {
			knownVersions = append(knownVersions, version)
			resolvedMethods[version] = make(map[string]*handler)
		}
	}

	// We must ensure that the earliest version is done first
	sort.Sort(knownVersions)

	var previousVersion string

	//	loop over versions, earliest first
	// build up each version by copying the previous version as a base
	// then setting on the version any explicitly defined method
	for _, version := range knownVersions {
		if previousVersion != "" {
			for mn, fn := range resolvedMethods[previousVersion] {
				resolvedMethods[version][mn] = fn
			}
		}

		for mn, fn := range s.registeredVersionMethods[version] {
			if fn == nil {
				delete(resolvedMethods[version], mn)
			} else {
				resolvedMethods[version][mn] = fn
			}
		}

		previousVersion = version
	}

	// build up latest methodSet if previous methodSet has been made
	if previousVersion != "" {
		resolvedMethods[VersionLatest] = resolvedMethods[previousVersion]
	}

	// Handle preview methods
	if len(s.registeredPreviewMethods) > 0 {
		resolvedMethods[VersionPreview] = make(map[string]*handler)
	}

	for mn, fn := range s.registeredPreviewMethods {
		if fn == nil {
			panic("cannot set preview method as nil")
		}

		resolvedMethods[VersionPreview][mn] = fn
	}

	s.resolvedMethods = resolvedMethods
}

// NoLongerSupported is a HandlerFunc which always returns the error
// `no_longer_supported` to the requester to indicate methods which have
// been withdrawn.
func NoLongerSupported(_ http.ResponseWriter, _ *Request) error {
	return cher.New(cher.NoLongerSupported, nil)
}

// Serve executes an RPC request.
func (s *Server) Serve(res http.ResponseWriter, req *Request) error {
	if s.AuthenticationMiddleware == nil {
		return cher.New(cher.AccessDenied, nil)
	}

	if s.resolvedMethods == nil {
		return cher.New("no_methods_registered", nil)
	}

	methodSet, ok := s.resolvedMethods[req.Version]
	if !ok {
		return cher.New(cher.NotFound, cher.M{"version": req.Version})
	}

	hn, ok := methodSet[req.Method]
	if !ok || hn == nil {
		return cher.New(cher.NotFound, cher.M{"method": req.Method, "version": req.Version})
	}

	// append latest version to Cuvva Endpoint Status
	appendCuvvaEndpointStatus(res, req.Version, hn.v)

	fn := hn.fn

	return fn(res, req)
}

const (
	// CuvvaEndpointStatus is the header appended to the response indicating the
	// usability status of the endpoint.
	CuvvaEndpointStatus = `Cuvva-Endpoint-Status`

	// PreviewNotice is the contents of the `Cuvva-Endpoint-Status` header when an
	// endpoint is called with the preview version.
	PreviewNotice = `preview; msg="endpoint is experimental and may change/be withdrawn without notice"`

	// LatestNotice is the contents of the `Cuvva-Endpoint-Status` header when an
	// endpoint is called to request the latest version.
	LatestNotice = `latest; msg="subject to change without notice"`

	// StableNotice is the contents of the `Cuvva-Endpoint-Status` header when an
	// endpoint is called and is not expected to change.
	StableNotice = `stable`
)

// appendCuvvaEndpointStatus applies the appropriate `Cuvva-Endpoint-Status`
// header for the method version requested by the client.
func appendCuvvaEndpointStatus(w http.ResponseWriter, requestedVersion, resolvedVersion string) {
	switch requestedVersion {
	case VersionPreview:
		w.Header().Set(CuvvaEndpointStatus, PreviewNotice)

	case VersionLatest:
		message := fmt.Sprintf(`%s; v="%s"`, LatestNotice, resolvedVersion)
		w.Header().Set(CuvvaEndpointStatus, message)

	default:
		w.Header().Set(CuvvaEndpointStatus, StableNotice)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.RawQuery != "" {
		s.writeError(w, cher.New("unexpected_input", nil))
		return
	}

	req := &Request{
		Body: r.Body,

		RemoteAddr:    r.RemoteAddr,
		BrowserOrigin: r.Header.Get("Origin"),
	}
	req.ctx = setRequestContext(r.Context(), req)

	var ok bool
	req.Method, req.Version, ok = requestPath(r.URL.Path)
	if !ok {
		s.writeError(w, cher.New(cher.NotFound, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	s.writeError(w, s.Serve(w, req))
}

// expRequestPath only matches HTTP Paths formed of /<version date>/<method name>
var expRequestPath = regexp.MustCompile(`^/(preview|latest|20\d{2}-\d{2}-\d{2})/([a-z0-9\_]+)$`)

func requestPath(path string) (method, version string, ok bool) {
	m := expRequestPath.FindStringSubmatch(path)
	if len(m) != 3 {
		return
	}

	version = m[1]
	method = m[2]
	ok = true
	return
}

func (s *Server) writeError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var body cher.E

	switch err := err.(type) {
	case cher.E:
		body = err

	case *json.SyntaxError:
		body = cher.New(
			"invalid_json",
			cher.M{
				"error":  err.Error(),
				"offset": err.Offset,
			},
		)

	case *json.UnmarshalTypeError:
		body = cher.New(
			"invalid_json",
			cher.M{
				"expected": err.Type.Kind().String(),
				"actual":   err.Value,
				"name":     err.Field,
			},
		)

	default:
		body = cher.New(cher.Unknown, nil)
	}

	w.WriteHeader(body.StatusCode())

	json.NewEncoder(w).Encode(body)
}
