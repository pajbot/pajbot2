package views

import (
	"html/template"
	"net/http"

	"github.com/pajlada/pajbot2/pkg/utils"
)

const (
	templatePrefix = "../../web/models/"
	templateSuffix = ".html"
	defaultTheme   = "default"
)

var cfg Config

type Config struct {
	WSHost string
}

type state struct {
	Config

	CurrentPage string

	LoggedIn bool

	Theme string
}

func Configure(c Config) {
	cfg = c
}

func templatePath(templateName string) string {
	return templatePrefix + templateName + templateSuffix
}

var validThemes = []string{"default", "dark"}

// get theme from cookie of request. return default if cookie is non-existant or invalid
func getTheme(r *http.Request) (theme string) {
	theme = defaultTheme

	themeCookie, err := r.Cookie("currentTheme")
	if err != nil {
		return
	}
	if themeCookie == nil {
		return
	}

	theme = themeCookie.Value

	if utils.StringContains(themeCookie.Value, validThemes) {
		theme = themeCookie.Value
	}

	return
}

func Render(templateName string, w http.ResponseWriter, r *http.Request) error {
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
		Theme:       getTheme(r),
	}

	err = tpl.ExecuteTemplate(w, "base"+templateSuffix, state)
	if err != nil {
		return err
	}

	return nil
}
