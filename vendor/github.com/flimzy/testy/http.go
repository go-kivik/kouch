package testy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

// ResponseHandler wraps an existing http.Response, to be served as a
// standard http.Handler
type ResponseHandler struct {
	*http.Response
}

var _ http.Handler = &ResponseHandler{}

func (h *ResponseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for header, values := range h.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	if h.StatusCode != 0 {
		w.WriteHeader(h.StatusCode)
	}
	if h.Body != nil {
		defer h.Body.Close() // nolint: errcheck
		io.Copy(w, h.Body)
	}
}

// ServeResponse starts a test HTTP server that serves r.
func ServeResponse(r *http.Response) *httptest.Server {
	return httptest.NewServer(&ResponseHandler{r})
}

// RequestValidator is a function that takes a *http.Request, and returns an
// error if it does not meet expectations. The error is turned into a 400
// response.
type RequestValidator func(*http.Request) error

// ValidateRequest returns a middleware that calls fn(), to validate the HTTP
// request, before continuing. An error returned by fn() will result in the
// addition of an X-Error header, a 400 status, and the error added to the
// body of the response.
func ValidateRequest(fn RequestValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := fn(r); err != nil {
				w.Header().Add("X-Error", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "%s", err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ServeResponseValidator wraps a ResponseHandler with ValidateRequest
// middleware for a complete response-serving, request-validating test server.
func ServeResponseValidator(r *http.Response, fn RequestValidator) *httptest.Server {
	mw := ValidateRequest(fn)
	return httptest.NewServer(mw(&ResponseHandler{r}))
}
