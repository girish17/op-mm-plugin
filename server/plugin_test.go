package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	assertVar := assert.New(t)
	plugin := Plugin{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	plugin.ServeHTTP(w, r)

	result := w.Result()
	assertVar.NotNil(result)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(result.Body)
	bodyBytes, err := io.ReadAll(result.Body)
	assertVar.Nil(err)
	bodyString := string(bodyBytes)

	assertVar.Equal("Hello, world!", bodyString)
}
