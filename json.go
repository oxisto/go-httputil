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
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// JSONResponseWithStatus returns a JSON encoded object with statusCode, if error is nil.
// Otherwise the error is returned and status code is set to http.StatusInternalServerError
func JSONResponseWithStatus(w http.ResponseWriter, r *http.Request, object interface{}, err error, statusCode int) {
	// uh-uh, we have an error
	if err != nil {
		logrus.Errorf("An error occured during processing of a REST request: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return not found if object is nil
	if object == nil {
		http.NotFound(w, r)
		return
	}

	// otherwise, lets try to decode the JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(object); err != nil {
		// uh-uh we couldn't decode the JSON
		logrus.Errorf("An error occured during encoding of the JSON response: %v. Payload was: %+v", err, object)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// JSONResponse returns a JSON encoded object with http.StatusOK, if error is nil.
// Otherwise the error is returned and status code is set to http.StatusInternalServerError
func JSONResponse(w http.ResponseWriter, r *http.Request, object interface{}, err error) {
	JSONResponseWithStatus(w, r, object, err, http.StatusOK)
}
