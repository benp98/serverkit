package assetkit

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

// ParseRootTemplates parses the files at the given http.Filesystem's root as HTML templates with their path as names
func ParseRootTemplates(fs http.FileSystem) (*template.Template, error) {
	dir, err := fs.Open("/")
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	items, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	itemNames := make([]string, 0)
	for _, item := range items {
		itemNames = append(itemNames, item.Name())
	}

	return ParseTemplates(fs, itemNames)
}

// ParseTemplates parses the given filepaths on the http.FileSystem as HTML templates with their path as template name
func ParseTemplates(fs http.FileSystem, items []string) (*template.Template, error) {
	tpl := template.New("")

	for _, item := range items {
		_, err := parseTemplateItem(tpl, item, fs)
		if err != nil {
			return nil, err
		}
	}

	return tpl, nil
}

// parseTemplateItem does all the dirty work. It takes the parent template, the filename and the http.FileSystem and reads the template file.
func parseTemplateItem(tpl *template.Template, name string, fs http.FileSystem) (*template.Template, error) {
	// Create new template
	t := tpl.New(name)

	// Open file for reading
	file, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file's content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Parse the template
	_, err = t.Parse(string(content))
	if err != nil {
		return nil, err
	}

	return t, err
}
