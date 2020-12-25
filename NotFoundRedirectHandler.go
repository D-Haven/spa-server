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
