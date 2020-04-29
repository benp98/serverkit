package httpkit

import (
	"net/http"
)

// The HandlerBuilder is used to set the default values for new instances of http.Handler
type HandlerBuilder struct {
	ErrorHandler ErrorHandler
}

// An ErrorHandler takes a http.ResponseWriter and any error values and writes its error message to the ResponseWriter
type ErrorHandler func(http.ResponseWriter, interface{})

// NewHandler returns a new generic HTTP handler
func (hb *HandlerBuilder) NewHandler(handleFunc func(*http.Request) (Response, error)) http.Handler {
	handler := new(genericHandler)

	handler.handleFunc = handleFunc

	// Set the error handler to the HandlerBuilder's value or fall back to the PlainTextErrorHandler if not defined
	if hb.ErrorHandler != nil {
		handler.errorHandler = hb.ErrorHandler
	} else {
		handler.errorHandler = PlainTextErrorHandler
	}

	return handler
}

type responseFunc func() Response
type genericHandler struct {
	handleFunc   func(*http.Request) (Response, error)
	errorHandler ErrorHandler
}

// ServeHTTP executes the real handler and processes its response or calls the error handler
func (h *genericHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var response Response
	var err error

	defer func() {
		// Normal error?
		if err != nil {
			h.errorHandler(res, err)
			return
		}

		// Panic?
		rec := recover()
		if rec != nil {
			h.errorHandler(res, rec)
			return
		}

		// No error -> write normal response
		err = response.WriteResponse(res)
		if err != nil {
			h.errorHandler(res, err)
		}
	}()

	response, err = h.handleFunc(req)
}
