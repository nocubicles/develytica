package utils

import (
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, filename string, data interface{}) error {
	tmpl, err := template.ParseGlob("static/*")

	if err != nil {
		return err
	}

	if err := tmpl.ExecuteTemplate(w, filename, data); err != nil {
		return err
	}

	return nil
}
