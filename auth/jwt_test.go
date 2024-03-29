/*
Copyright 2019 Christian Banse

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestHandleWithNext(t *testing.T) {
	token, err := IssueToken([]byte("secret"), "me", time.Now().Add(time.Hour*24))

	t.Logf(token.AccessToken)

	assert.Nil(t, err)

	r := http.Request{
		Header: http.Header{
			"Authorization": []string{"Bearer " + token.AccessToken},
		},
	}

	options := DefaultOptions
	options.JWTKeySupplier = func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	}
	options.TokenExtractor = ExtractFromFirstAvailable(
		ExtractTokenFromCookie("auth"),
		ExtractTokenFromHeader)

	handler := NewHandler(options)

	c := &gin.Context{Request: &r, KeysMutex: &sync.RWMutex{}}
	handler.AuthRequired(c)

	parsed := c.Value(ClaimsContext)

	assert.NotNil(t, parsed)
}

func TestExctractTokenFromHeader(t *testing.T) {
	r := http.Request{
		Header: http.Header{
			"Authorization": []string{"Bearer token"},
		},
	}

	token, err := ExtractTokenFromHeader(&r)

	assert.Nil(t, err)
	assert.Equal(t, "token", token)
}
