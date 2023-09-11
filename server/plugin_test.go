package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	assertVar := assert.New(t)
	p := Plugin{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	c := plugin.Context{
		SessionId:      "",
		RequestId:      "",
		IPAddress:      "",
		AcceptLanguage: "",
		UserAgent:      "",
	}
	p.ServeHTTP(&c, w, r)

	result := w.Result()
	assertVar.NotNil(result)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(result.Body)
	bodyBytes, err := io.ReadAll(result.Body)
	assertVar.Nil(err)
	bodyString := string(bodyBytes)

	assertVar.Equal("404 page not found\n", bodyString)
}
