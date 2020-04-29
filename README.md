# ServerKit

ServerKit is a collection of utility functions for building HTTP servers with a nicer API in Go. It started as part of a private project which was carried over to other private projects as well. This is a cleaned-up version of that code.

## Features

* Nicer API for HTTP Handlers -> Handler function takes request and returns a response struct
* Uses gzip compression if enabled in the handler function and supported by the client
* Loading of HTML templates from http.FileSystem, useful for template files packed into the executable using tools like [fileb0x](https://github.com/UnnoTed/fileb0x)

## How it looks
```go
handlerBuilder := &httpkit.HandlerBuilder{
	ErrorHandler: httpkit.SimpleHTMLErrorHandler,
}

mux := http.NewServeMux()

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
```

## Error pages

The ServerKit handler displays an error page if the handler function returns an error value or panics. To display custom error pages or enable extended logging of the error, a custom error handler can be provided.

## API Stability
The stability of this module's API is not guaranteed and features may be added or evaluated and removed at any point.
