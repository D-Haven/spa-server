/*
 * Copyright (c) 2020. D-Haven
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"log"
	"net/http"
)

func NotFoundRedirectHandler(redirectTarget string, handler http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		wrapper := &NotFoundRedirectWrapper{
			ResponseWriter: writer,
			path:           request.URL.Path,
		}

		handler.ServeHTTP(wrapper, request)

		if wrapper.status == http.StatusNotFound && wrapper.path != "/favicon.ico" {
			log.Printf("Redirecting %s to %s.", request.RequestURI, redirectTarget)
			http.Redirect(writer, request, redirectTarget, http.StatusSeeOther)
		}
	}
}

type NotFoundRedirectWrapper struct {
	http.ResponseWriter // We embed http.ResponseWriter
	status              int
	path                string
}

func (wrapper *NotFoundRedirectWrapper) WriteHeader(status int) {
	wrapper.status = status

	if status != http.StatusNotFound || wrapper.path == "/favicon.ico" {
		wrapper.ResponseWriter.WriteHeader(status)
	}
}

func (wrapper *NotFoundRedirectWrapper) Write(content []byte) (int, error) {
	if wrapper.status != http.StatusNotFound || wrapper.path == "/favicon.ico" {
		return wrapper.ResponseWriter.Write(content)
	}

	// Lie that we successfully written it
	return len(content), nil
}
