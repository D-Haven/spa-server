/*
 * Copyright (c) 2021. D-Haven
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

package handlers

import (
	"log"
	"net/http"
)

// Browsers will pull favicon.ico whether you have a link to it or not
// Apple devices will pull site.webmanifest
// Search engines look for robots.txt
// All of these paths should return 404 if the file does not exist rather than a redirect
var protectedPaths = []string{
	"/favicon.ico",
	"/robots.txt",
	"/site.webmanifest",
	"/sitemap.xml",
	"/search.xml",
}

type NotFoundRedirectWrapper struct {
	http.ResponseWriter // We embed http.ResponseWriter
	status              int
	path                string
}

func isIn(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}

func NotFoundRedirectHandler(handler http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		wrapper := &NotFoundRedirectWrapper{
			ResponseWriter: writer,
			path:           request.URL.Path,
		}

		handler.ServeHTTP(wrapper, request)

		if wrapper.status == http.StatusNotFound && !isIn(wrapper.path, protectedPaths) {
			log.Printf("Redirecting %s to /.", request.RequestURI)
			http.Redirect(writer, request, "/", http.StatusSeeOther)
		}
	}
}

func (wrapper *NotFoundRedirectWrapper) WriteHeader(status int) {
	wrapper.status = status

	if status != http.StatusNotFound || isIn(wrapper.path, protectedPaths) {
		wrapper.ResponseWriter.WriteHeader(status)
	}
}

func (wrapper *NotFoundRedirectWrapper) Write(content []byte) (int, error) {
	if wrapper.status != http.StatusNotFound || isIn(wrapper.path, protectedPaths) {
		return wrapper.ResponseWriter.Write(content)
	}

	// Lie that we successfully written it
	return len(content), nil
}
