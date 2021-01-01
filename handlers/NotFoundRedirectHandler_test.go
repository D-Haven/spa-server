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
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TestWrapperRedirectsToRootIfNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	CheckError(err)

	rr := httptest.NewRecorder()

	handler := NotFoundRedirectHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err = io.WriteString(w, "{}")
		CheckError(err)
	}))

	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusSeeOther {
		t.Errorf("Redirect was supposed to happen, but it was %d", rr.Result().StatusCode)
	}

	if rr.Result().Header.Get("Location") != "/" {
		t.Errorf("Redirect must be '/', but was '%s'", rr.Result().Header.Get("Location"))
	}
}

func TestWrapperPassesFoundPathsUnharmed(t *testing.T) {
	req, err := http.NewRequest("GET", "/test.txt", nil)
	CheckError(err)

	content := "This is a test"
	rr := httptest.NewRecorder()

	handler := NotFoundRedirectHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err = io.WriteString(w, content)
		CheckError(err)
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected HTTP OK, but was %d", rr.Code)
	}

	if rr.Body.String() != content {
		t.Errorf("Expected original content, but was:\n%s", rr.Body.String())
	}
}

func TestWrapperAllowsNotFoundForSpecialFiles(t *testing.T) {
	protectedPaths = []string{
		"/favicon.ico",
		"/robots.txt",
		"/site.webmanifest",
		"/sitemap.xml",
		"/search.xml",
	}

	for _, path := range protectedPaths {
		t.Run(path, func(t *testing.T) {
			req, err := http.NewRequest("GET", path, nil)
			CheckError(err)

			rr := httptest.NewRecorder()

			handler := NotFoundRedirectHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_, err = io.WriteString(w, "{}")
				CheckError(err)
			}))

			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusNotFound {
				t.Errorf("Expected HTTP Not Found, but was %d", rr.Code)
			}

			if rr.Body.String() != "{}" {
				t.Errorf("Body was changed:\n%s", rr.Body.String())
			}
		})
	}
}
