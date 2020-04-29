package httpkit

// This file contains some predefined Handlers and Responses

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

// Handle404Plain returns a plain text answer for the 404 error
func Handle404Plain(req *http.Request) Response {
	res := NewPlainTextResponse(404)
	fmt.Fprintln(res, "Error 404: Not found")

	return res
}

// Handle404SimpleHTML returns a simple HTML answer for the 404 error
func Handle404SimpleHTML(req *http.Request) Response {
	res := NewHTMLResponse(404)
	fmt.Fprintln(res, simpleHTMLTemplate("Error 404: Not Found", "The resource was not found on this server."))

	return res
}

// NewPlainTextResponse returns a new Plaintext response with the given status code
func NewPlainTextResponse(status int) ExtendedResponse {
	return NewResponse(status, "text/plain", true)
}

// NewHTMLResponse returns a new HTML response with the given status code
func NewHTMLResponse(status int) ExtendedResponse {
	return NewResponse(status, "text/html", true)
}

// NewJSONResponse returns a new JSON response with the given status code and content
func NewJSONResponse(status int, data interface{}) ExtendedResponse {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	res := NewResponse(status, "application/json", true)

	_, err = res.Write(content)
	if err != nil {
		panic(err)
	}

	return res
}

// NewRedirectResponse returns a new Redirect Response
func NewRedirectResponse(permanent bool, location string) ExtendedResponse {
	var status int
	if permanent {
		status = 301
	} else {
		status = 302
	}

	res := NewResponse(status, "text/html", true)
	res.SetHeader("Location", location)
	fmt.Fprintln(res, simpleHTMLTemplate("Redirect", fmt.Sprintf(
		`You see this Text because your browser cannot handle this redirect. Please click the following link to go to the destination: <a href="%s">Go to destination</a>`,
		location,
	)))

	return res
}

// PlainTextErrorHandler writes the error message as plain text to the HTTP response
func PlainTextErrorHandler(res http.ResponseWriter, err interface{}) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(500)
	fmt.Fprintln(res, "Error:")
	fmt.Fprintln(res, err)
}

// SimpleHTMLErrorHandler writes the error message as simple HTML to the HTTP response
func SimpleHTMLErrorHandler(res http.ResponseWriter, err interface{}) {
	t := template.Must(template.New("").Parse(`<!DOCTYPE html><html><head><title>Error</title><meta charset="utf-8"></head><body><h1>Error</h1><p>{{.}}</p></body></html>`))
	res.Header().Add("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(500)
	t.Execute(res, err)
}
