package views

import (
	"html/template"
	"net/http"
)

const templatePrefix = "../../web/models/"
const templateSuffix = ".html"

var cfg Config

type Config struct {
	WSHost string
}

type state struct {
	Config

	CurrentPage string

	LoggedIn bool
}

func Configure(c Config) {
	cfg = c
}

func templatePath(templateName string) string {
	return templatePrefix + templateName + templateSuffix
}

func Render(templateName string, w http.ResponseWriter) error {
	var err error
	tpl := template.New(templateName)
	_, err = tpl.ParseFiles(templatePath(templateName), templatePath("base"))
	if err != nil {
		return err
	}

	state := state{
		Config: cfg,

		CurrentPage: templateName,
		LoggedIn:    false,
	}

	err = tpl.ExecuteTemplate(w, "base"+templateSuffix, state)
	if err != nil {
		return err
	}

	return nil
}
