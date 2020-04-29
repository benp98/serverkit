package httpkit

import "fmt"

// simpleHTMLTemplate creates a simple HTML message
func simpleHTMLTemplate(title string, content interface{}) string {
	return fmt.Sprintf(`<!DOCTYPE html><html><head><title>%s</title><meta charset="utf-8"></head><body><h1>%s</h1><p>%s</p></body></html>`, title, title, content)
}
