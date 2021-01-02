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
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

type GzipResponseWrapper struct {
	io.Writer
	http.ResponseWriter
}

func (wrapper GzipResponseWrapper) Write(content []byte) (int, error) {
	if "" == wrapper.Header().Get("Content-Type") {
		// If no content type, apply sniffing algorithm to un-gzipped body.
		wrapper.Header().Set("Content-Type", http.DetectContentType(content))
	}

	return wrapper.Writer.Write(content)
}

func GzipHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			handler(writer, request)
			return
		}

		writer.Header().Set("Content-Encoding", "gzip")
		compressor := gzip.NewWriter(writer)
		defer func() {
			if err := compressor.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		gzipWrapper := GzipResponseWrapper{Writer: compressor, ResponseWriter: writer}
		handler(gzipWrapper, request)
	}
}