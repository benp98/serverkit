package httpkit

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Response represents a (before calling WriteResponse non-sent) HTTP response returned by the handler function.
type Response interface {
	WriteResponse(http.ResponseWriter) error
}

// ExtendedResponse is a extension of the Response interface which adds some comfort functions and extended Features
type ExtendedResponse interface {
	Response
	io.Writer // Append data to the content

	SetHeader(string, string)
	DisableCaching()
	AllowCompression(*http.Request)
}

// The actual implementation of Response
type genericResponse struct {
	status      int
	contentType string
	isText      bool
	content     bytes.Buffer
	headers     map[string]string
	compress    bool
}

// NewResponse returns a new Response with the given values
func NewResponse(status int, contentType string, isText bool) ExtendedResponse {
	response := new(genericResponse)
	response.status = status
	response.contentType = contentType
	response.isText = isText
	response.headers = make(map[string]string)
	return response
}

// Write appends the given data to the contents of this resource
func (res *genericResponse) Write(data []byte) (int, error) {
	return res.content.Write(data)
}

// WriteResponse writes the HTTP response to the given http.ResponseWriter
func (res *genericResponse) WriteResponse(rw http.ResponseWriter) error {
	// Set Content Type Headers
	if res.isText {
		rw.Header().Add("Content-Type", fmt.Sprintf("%s; charset=utf-8", res.contentType))
	} else {
		rw.Header().Add("Content-Type", res.contentType)
	}

	if res.compress {
		rw.Header().Add("Content-Encoding", "gzip")
	}

	// Set Custom Headers
	for k, v := range res.headers {
		rw.Header().Set(k, v)
	}

	// Write the Response
	rw.WriteHeader(res.status)

	// Handle compression
	var err error
	if res.compress {
		gz := gzip.NewWriter(rw)
		_, err = gz.Write(res.content.Bytes())
		if err != nil {
			return err
		}
		err = gz.Close()
	} else {
		_, err = rw.Write(res.content.Bytes())
	}

	return err
}

// SetHeader sets the given custom header to the provided value
func (res *genericResponse) SetHeader(key, value string) {
	res.headers[key] = value
}

// DisableCaching sets the Cache-Control, Pragma and Expires HTTP Headers to values which tell the browser not to cache the response.
func (res *genericResponse) DisableCaching() {
	res.SetHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	res.SetHeader("Pragma", "no-cache")
	res.SetHeader("Expires", "0")
}

func (res *genericResponse) AllowCompression(req *http.Request) {
	res.compress = strings.Contains(req.Header.Get("Accept-Encoding"), "gzip")
}
