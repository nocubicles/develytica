package utils

import (
	"fmt"
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, filename string, data interface{}) error {
	tmpl, err := template.ParseGlob("templates/*")

	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := tmpl.ExecuteTemplate(w, filename, data); err != nil {
		fmt.Println(err)
	}

	return nil
}
