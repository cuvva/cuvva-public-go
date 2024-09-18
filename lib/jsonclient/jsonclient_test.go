package jsonclient

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestGetHTTPMethod(t *testing.T) {
	defer gock.Off()

	gock.New("http://coo.va/").
		Get("/test").
		Reply(http.StatusNoContent)

	client := NewClient("http://coo.va/", nil)
	gock.InterceptClient(client.Client)

	err := client.Do(context.Background(), "GET", "test", nil, nil, nil)

	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func TestPutHTTPMethod(t *testing.T) {
	defer gock.Off()

	gock.New("http://coo.va/").
		Put("/test").
		Reply(http.StatusNoContent)

	client := NewClient("http://coo.va/", nil)
	gock.InterceptClient(client.Client)

	err := client.Do(context.Background(), "PUT", "test", nil, nil, nil)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func TestPostHTTPMethod(t *testing.T) {
	defer gock.Off()

	gock.New("http://coo.va/").
		Post("/test").
		Reply(http.StatusNoContent)

	client := NewClient("http://coo.va/", nil)
	gock.InterceptClient(client.Client)

	err := client.Do(context.Background(), "POST", "test", nil, nil, nil)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func TestDeleteHTTPMethod(t *testing.T) {
	defer gock.Off()

	gock.New("http://coo.va/").
		Delete("/test").
		Reply(http.StatusNoContent)

	client := NewClient("http://coo.va/", nil)
	gock.InterceptClient(client.Client)

	err := client.Do(context.Background(), "DELETE", "test", nil, nil, nil)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func TestRequestQuery(t *testing.T) {
	defer gock.Off()

	paramKey := "testing"
	paramValue := "true"

	gock.New("http://coo.va/").
		Get("/test").
		MatchParam(paramKey, paramValue).
		Reply(http.StatusNoContent)

	client := NewClient("http://coo.va/", nil)
	gock.InterceptClient(client.Client)

	err := client.Do(context.Background(), "GET", "test", url.Values{paramKey: {paramValue}}, nil, nil)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func TestRequestBody(t *testing.T) {
	defer gock.Off()

	testJSON := map[string]bool{"testing": true}

	gock.New("http://coo.va/").
		Post("/test").
		MatchType("application/json; charset=utf-8").
		JSON(testJSON).
		Reply(http.StatusNoContent)

	client := NewClient("http://coo.va/", nil)
	gock.InterceptClient(client.Client)

	err := client.Do(context.Background(), "POST", "test", nil, testJSON, nil)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func TestRequestModifier(t *testing.T) {
	defer gock.Off()

	testJSON := map[string]bool{"testing": true}

	modifier := func(req *http.Request) {
		req.Header.Add("X-Test-Header", "test")
	}

	gock.New("http://coo.va/").
		Post("/test").
		MatchType("application/json; charset=utf-8").
		JSON(testJSON).
		MatchHeader("X-Test-Header", "test").
		Reply(http.StatusNoContent)

	client := NewClient("http://coo.va/", nil)
	gock.InterceptClient(client.Client)

	err := client.Do(context.Background(), "POST", "test", nil, testJSON, nil, modifier)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func TestResponseBody(t *testing.T) {
	defer gock.Off()

	gock.New("http://coo.va/").
		Get("/test").
		MatchHeader("Accept", "application/json").
		Reply(http.StatusOK).
		JSON(map[string]bool{"testing": true})

	client := NewClient("http://coo.va/", nil)
	gock.InterceptClient(client.Client)

	var response map[string]bool
	err := client.Do(context.Background(), "GET", "test", nil, nil, &response)
	assert.Nil(t, err)
	assert.True(t, response["testing"])
	assert.True(t, gock.IsDone())
}

func TestErrorUnmarshaling(t *testing.T) {
	defer gock.Off()

	responseError := cher.E{Code: "test_error"}

	gock.New("http://coo.va/").
		Get("/test").
		Reply(http.StatusBadRequest).
		JSON(responseError)

	client := NewClient("http://coo.va/", nil)
	gock.InterceptClient(client.Client)

	err := client.Do(context.Background(), "GET", "test", nil, nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, responseError.Code, err.(cher.E).Code)
	assert.True(t, gock.IsDone())
}

func TestErrorCatching(t *testing.T) {
	defer gock.Off()

	gock.New("http://coo.va/").
		Get("/test").
		Reply(http.StatusInternalServerError)

	client := NewClient("http://coo.va/", nil)
	gock.InterceptClient(client.Client)

	err := client.Do(context.Background(), "GET", "test", nil, nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "internal_server_error", err.(cher.E).Code)
	assert.True(t, gock.IsDone())
}
