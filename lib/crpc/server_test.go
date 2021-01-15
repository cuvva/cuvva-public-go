//nolint:bodyclose // incorrect
package crpc

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
)

type testInput struct{}

type testOutput struct{}

func TestWrap(t *testing.T) {
	tests := []struct {
		Name  string
		Fn    interface{}
		Error error
	}{
		{
			"HandlerFunc",
			HandlerFunc(func(w http.ResponseWriter, r *Request) error { return nil }),
			errors.New("fn doesn't need to be wrapped, use RegisterFunc"),
		},
		{
			"WrappedFunc", &WrappedFunc{},
			errors.New("fn is already wrapped, use RegisterFunc"),
		},
		{
			"NotFunc", "string",
			errors.New("fn must be function, got string"),
		},
		{
			"NoInput", func() {},
			errors.New("fn input must be (context.Context) or (context.Context, *T), got 0 arguments"),
		},
		{
			"LongInput", func(ctx context.Context, foo string, bar string) {},
			errors.New("fn input must be (context.Context) or (context.Context, *T), got 3 arguments"),
		},
		{
			"NoOutput", func(ctx context.Context) {},
			errors.New("fn output must be (error) or (*T, error), got 0 arguments"),
		},
		{
			"LongOutput", func(ctx context.Context) (foo, bar string, err error) { return },
			errors.New("fn output must be (error) or (*T, error), got 3 arguments"),
		},
		{
			"ContextRequired", func(foo string) error { return nil },
			errors.New("fn first argument must implement context.Context, got string"),
		},
		{
			"ErrorRequired", func(ctx context.Context) string { return "" },
			errors.New("fn last argument must implement error, got string"),
		},
		{
			"InputNotPointer", func(ctx context.Context, in testInput) error { return nil },
			errors.New("fn last argument must be a pointer, got crpc.testInput"),
		},
		{
			"InputNotStruct", func(ctx context.Context, in *string) error { return nil },
			errors.New("fn last argument must be a struct, got string"),
		},
		{
			"OutputNotPointer", func(ctx context.Context) (out testOutput, err error) { return },
			errors.New("unsupported return type, expected *struct or slice; got crpc.testOutput"),
		},
		{
			"OutputNotStructSlice", func(ctx context.Context) (out *string, err error) { return },
			errors.New("unsupported return type, expected *struct or slice; got string"),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			_, err := Wrap(test.Fn)
			if test.Error != nil {
				assert.Equal(t, test.Error, err)
			}
		})
	}
}

func TestMethodsAreBroughtForward(t *testing.T) {
	foov1 := &handler{v: "2019-01-01"}
	barv1 := &handler{v: "2019-02-02"}

	zs := Server{
		registeredVersionMethods: map[string]map[string]*handler{
			"2019-01-01": {
				"foo": foov1,
			},
			"2019-02-02": {
				"bar": barv1,
			},
		},
	}

	zs.buildRoutes()

	expected := map[string]map[string]*handler{
		"2019-01-01": {
			"foo": foov1,
		},
		"2019-02-02": {
			"foo": foov1,
			"bar": barv1,
		},
		"latest": {
			"foo": foov1,
			"bar": barv1,
		},
	}

	assert.Equal(t, expected, zs.resolvedMethods)
}

func TestMethodsAreBroughtForwardComplex(t *testing.T) {
	foov1 := &handler{v: "2019-01-01"}
	foov2 := &handler{v: "2019-02-02"}
	barv1 := &handler{v: "2019-02-02"}
	barv2 := &handler{v: "2019-03-03"}
	barv3 := &handler{v: "2019-04-04"}

	zs := Server{
		registeredVersionMethods: map[string]map[string]*handler{
			"2019-01-01": {
				"foo": foov1,
			},
			"2019-02-02": {
				"foo": foov2,
				"bar": barv1,
			},
			"2019-03-03": {
				"bar": barv2,
			},
			"2019-04-04": {
				"bar": barv3,
			},
		},
	}

	zs.buildRoutes()

	expected := map[string]map[string]*handler{
		"2019-01-01": {
			"foo": foov1,
		},
		"2019-02-02": {
			"foo": foov2,
			"bar": barv1,
		},
		"2019-03-03": {
			"foo": foov2,
			"bar": barv2,
		},
		"2019-04-04": {
			"foo": foov2,
			"bar": barv3,
		},
		"latest": {
			"foo": foov2,
			"bar": barv3,
		},
	}

	assert.Equal(t, expected, zs.resolvedMethods)
}

func TestMethodsAreBroughtForwardAndRemoved(t *testing.T) {
	foov1 := &handler{v: "2019-01-01"}
	barv1 := &handler{v: "2019-01-01"}

	zs := Server{
		registeredVersionMethods: map[string]map[string]*handler{
			"2019-01-01": {
				"foo": foov1,
				"bar": barv1,
			},
			"2019-02-02": {
				"bar": nil,
			},
		},
	}

	zs.buildRoutes()

	expected := map[string]map[string]*handler{
		"2019-01-01": {
			"foo": foov1,
			"bar": barv1,
		},
		"2019-02-02": {
			"foo": foov1,
		},
		"latest": {
			"foo": foov1,
		},
	}

	assert.Equal(t, expected, zs.resolvedMethods)
}

func TestMethodsAreDefinedRemovedMultiple(t *testing.T) {
	foov1 := &handler{v: "2019-01-01"}
	barv1 := &handler{v: "2019-01-01"}
	foov2 := &handler{v: "2019-02-02"}
	foov3 := &handler{v: "2019-03-03"}

	zs := Server{
		registeredVersionMethods: map[string]map[string]*handler{
			"2019-01-01": {
				"foo": foov1,
				"bar": barv1,
			},
			"2019-02-02": {
				"foo": nil,
			},
			"2019-03-03": {
				"foo": foov2,
			},
			"2019-04-04": {
				"foo": nil,
			},
			"2019-05-05": {
				"foo": foov3,
			},
		},
	}

	zs.buildRoutes()

	expected := map[string]map[string]*handler{
		"2019-01-01": {
			"foo": foov1,
			"bar": barv1,
		},
		"2019-02-02": {
			"bar": barv1,
		},
		"2019-03-03": {
			"foo": foov2,
			"bar": barv1,
		},
		"2019-04-04": {
			"bar": barv1,
		},
		"2019-05-05": {
			"foo": foov3,
			"bar": barv1,
		},
		"latest": {
			"foo": foov3,
			"bar": barv1,
		},
	}

	assert.Equal(t, expected, zs.resolvedMethods)
}

func TestPreviewMethodsAreRegistered(t *testing.T) {
	barv1 := &handler{v: "2019-01-01"}
	fooPrev := &handler{v: "preview"}

	zs := Server{
		registeredPreviewMethods: map[string]*handler{
			"foo": fooPrev,
		},
		registeredVersionMethods: map[string]map[string]*handler{
			"2019-01-01": {
				"bar": barv1,
			},
		},
	}

	zs.buildRoutes()

	expected := map[string]map[string]*handler{
		"preview": {
			"foo": fooPrev,
		},
		"2019-01-01": {
			"bar": barv1,
		},
		"latest": {
			"bar": barv1,
		},
	}

	assert.Equal(t, expected, zs.resolvedMethods)
}

func TestNilPreviewMethodsPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should panic when setting preview methods to nil")
		}
	}()

	zs := NewServer(UnsafeNoAuthentication)

	zs.Register("foo", "preview", nil, nil)
}

func TestPanicIfMethodDeclaredTwice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should panic when declaring same method, same version, twice")
		}
	}()

	zs := NewServer(UnsafeNoAuthentication)

	zs.Register("foo", "2019-01-01", nil, func(context.Context) error { return nil })

	zs.Register("foo", "2019-01-01", nil, func(context.Context) error { return nil })
}

func TestMiddlewareIsLoadedInOrder(t *testing.T) {
	zs := NewServer(UnsafeNoAuthentication)

	zs.Register("foo", "preview", nil, makeRPCCall("called foo!"))
	zs.Use(addHeaderMiddleware("X-Is-Test", "win!"))
	zs.Register("bar", "preview", nil, makeRPCCall("called bar!"))

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/preview/foo", nil)

	zs.ServeHTTP(w, r)

	if _, ok := w.Result().Header["X-Is-Test"]; ok {
		t.Error("was expecting 'X-Is-Test' header to not be present")
	}

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("POST", "/preview/bar", nil)

	zs.ServeHTTP(w, r)

	if _, ok := w.Result().Header["X-Is-Test"]; !ok {
		t.Error("was expecting 'X-Is-Test' header to be present")
	}
}

func TestMiddlewareRunsGlobalInOrderAndRequestSpecific(t *testing.T) {
	zs := NewServer(UnsafeNoAuthentication)

	zs.Use(addHeaderMiddleware("X-Present-On-Both", "win!"))
	zs.Register("foo", "preview", nil, makeRPCCall("called foo!"))
	zs.Use(addHeaderMiddleware("X-Present-On-Bar", "win!"))
	zs.Register("bar", "preview", nil, makeRPCCall("called bar!"), addHeaderMiddleware("X-Also-On-Bar", "wat?"))

	w1 := httptest.NewRecorder()
	w2 := httptest.NewRecorder()
	r1, _ := http.NewRequest("POST", "/preview/foo", nil)
	r2, _ := http.NewRequest("POST", "/preview/bar", nil)

	zs.ServeHTTP(w1, r1)
	zs.ServeHTTP(w2, r2)

	if _, ok := w1.Result().Header["X-Present-On-Both"]; !ok {
		t.Error("was expecting 'X-Present-On-Both' header to be present")
	}

	if _, ok := w2.Result().Header["X-Present-On-Both"]; !ok {
		t.Error("was expecting 'X-Present-On-Both' header to be present")
	}

	if _, ok := w1.Result().Header["X-Present-On-Bar"]; ok {
		t.Error("was expecting 'X-Present-On-Bar' header to NOT be present")
	}

	if _, ok := w2.Result().Header["X-Present-On-Bar"]; !ok {
		t.Error("was expecting 'X-Present-On-Bar' header to be present")
	}

	if _, ok := w1.Result().Header["X-Also-On-Bar"]; ok {
		t.Error("was expecting 'X-Also-On-Bar' header to NOT be present")
	}

	if _, ok := w2.Result().Header["X-Also-On-Bar"]; !ok {
		t.Error("was expecting 'X-Also-On-Bar' header to be present")
	}
}

type testResponse struct {
	Message string `json:"message"`
}

func addHeaderMiddleware(headerToAdd, value string) func(HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(res http.ResponseWriter, req *Request) error {
			res.Header().Add(headerToAdd, value)

			return next(res, req)
		}
	}
}

func makeRPCCall(messageToReturn string) func(context.Context) (*testResponse, error) {
	return func(_ context.Context) (*testResponse, error) {
		return &testResponse{
			Message: messageToReturn,
		}, nil
	}
}

func TestSchemasAreCompiled(t *testing.T) {
	brokenSchema := gojsonschema.NewStringLoader(`{
		"type": "object",
		"properties":}
	}`)
	validSchema := gojsonschema.NewStringLoader(`{
		"type": "object",
		"properties": {
			"foo": {
				"type": "string"
			}
		}
	}`)

	handler := func(_ context.Context, _ *struct{}) error {
		return nil
	}

	zs := NewServer(UnsafeNoAuthentication)

	zs.Register("should_pass", "2019-01-01", validSchema, handler)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("should_crash method should panic")
		}
	}()

	zs.Register("should_crash", "2019-01-01", brokenSchema, handler)
}

func UnsafeNoAuthentication(next HandlerFunc) HandlerFunc {
	return func(res http.ResponseWriter, req *Request) error {
		return next(res, req)
	}
}
