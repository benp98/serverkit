package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/benp98/serverkit/assetkit"
	"github.com/benp98/serverkit/httpkit"
)

// This is a sample web server application which utilises the functions of ServerKit
func main() {
	server := &http.Server{}
	server.Addr = "localhost:8080"

	log.Print("Parsing templates")
	templateFS := http.Dir("template")
	templates, err := assetkit.ParseRootTemplates(templateFS)
	if err != nil {
		log.Fatal(err)
	}

	handlerBuilder := &httpkit.HandlerBuilder{
		ErrorHandler: httpkit.SimpleHTMLErrorHandler,
	}

	log.Print("Setting up HTTP handlers")
	mux := http.NewServeMux()
	server.Handler = mux

	// The root handler displays a template based HTML page
	mux.Handle("/", handlerBuilder.NewHandler(func(req *http.Request) (httpkit.Response, error) {
		// Filter out any requests not for "/" and display a 404 page
		if req.URL.Path != "/" {
			return httpkit.Handle404SimpleHTML(req), nil
		}

		res := httpkit.NewHTMLResponse(200)
		res.AllowCompression(req) // Allow compression of the response if req contains a valid header
		err := templates.ExecuteTemplate(res, "index.html", []string{"This is an example", "for github.com/benp98/serverkit"})
		return res, err // We don't need an error check, because the the real handler function checks if err is not nil
	}))

	// A simple redirect
	mux.Handle("/redirect", handlerBuilder.NewHandler(func(req *http.Request) (httpkit.Response, error) {
		return httpkit.NewRedirectResponse(false, "/"), nil
	}))

	// JSON example
	message := "Hello"
	mux.Handle("/api", handlerBuilder.NewHandler(func(req *http.Request) (httpkit.Response, error) {
		// Behaviour depends on HTTP method
		switch req.Method {
		case http.MethodPut:
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}

			message = string(body)

			return httpkit.NewJSONResponse(200, struct {
				Updated bool
			}{
				true,
			}), nil
		default:
			return httpkit.NewJSONResponse(200, struct {
				Message string
			}{
				message,
			}), nil
		}
	}))

	log.Printf("Listening on %s", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
