package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (writer GzipResponseWriter) Write(content []byte) (int, error) {
	if "" == writer.Header().Get("Content-Type") {
		// If no content type, apply sniffing algorithm to un-gzipped body.
		writer.Header().Set("Content-Type", http.DetectContentType(content))
	}

	return writer.Writer.Write(content)
}

func GzipHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			handler(writer, request)
			return
		}

		writer.Header().Set("Content-Encoding", "gzip")
		compressor := gzip.NewWriter(writer)
		defer CheckError(compressor.Close())

		gzipWriter := GzipResponseWriter{Writer: compressor, ResponseWriter: writer}
		handler(gzipWriter, request)
	}
}
