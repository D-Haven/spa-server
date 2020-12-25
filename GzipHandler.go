package main

import (
	"compress/gzip"
	"io"
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
		defer func() { CheckError(compressor.Close()) }()

		gzipWrapper := GzipResponseWrapper{Writer: compressor, ResponseWriter: writer}
		handler(gzipWrapper, request)
	}
}
