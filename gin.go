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

package httputil

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func IntParam(c *gin.Context, key string) (i int64, err error) {
	return strconv.ParseInt(c.Param(key), 10, 64)
}

func FloatParam(c *gin.Context, key string) (i float64, err error) {
	return strconv.ParseFloat(c.Param(key), 64)
}

func IntQuery(c *gin.Context, key string) (i int64, err error) {
	return strconv.ParseInt(c.Query(key), 10, 64)
}

func FloatQuery(c *gin.Context, key string) (i float64, err error) {
	return strconv.ParseFloat(c.Query(key), 64)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSON(c *gin.Context, status int, value interface{}, err error) {
	if err != nil {
		logrus.Errorf("An error occurred during processing of a REST request: %s", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		c.JSON(http.StatusNotFound, nil)
	} else {
		c.JSON(status, value)
	}
}
