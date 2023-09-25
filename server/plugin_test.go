/*
 *    Copyright 2023 Girish M
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

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
