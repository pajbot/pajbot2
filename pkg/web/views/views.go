package views

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/pajbot/pajbot2/internal/config"
	"github.com/pajbot/utils"
)

const (
	templateSuffix = ".html"
	defaultTheme   = "Dark"
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

	Extra string
}

func Configure(c Config) {
	cfg = c
}

func templatePath(templateName string) string {
	return filepath.Join(config.WebStaticPath, "views", templateName+templateSuffix)
}

var validThemes = []string{"Light", "Dark"}

// get theme from cookie of request. return default if cookie is non-existent or invalid
func getTheme(r *http.Request) (theme string) {
	theme = defaultTheme

	themeCookie, err := r.Cookie("currentTheme")
	if err != nil || themeCookie == nil {
		// no cookie to be found
		return
	}

	// verify that theme name in cookie is valid
	if utils.StringContains(themeCookie.Value, validThemes) {
		theme = themeCookie.Value
	}

	return
}

// RenderBasic renders only the given template, no base files
func RenderBasic(templateName string, w http.ResponseWriter, r *http.Request) error {
	tpl, err := template.ParseFiles(templatePath(templateName))
	if err != nil {
		return err
	}

	state := state{
		Config: cfg,
		Theme:  getTheme(r),
	}

	err = tpl.Execute(w, state)
	if err != nil {
		return err
	}

	return nil
}

func RenderExtra(templateName string, w http.ResponseWriter, r *http.Request, extra []byte) error {
	var err error
	tpl := template.New(templateName)
	_, err = tpl.ParseFiles(templatePath(templateName), templatePath("index"))
	if err != nil {
		return err
	}

	state := state{
		Config: cfg,

		CurrentPage: templateName,
		LoggedIn:    false,
		Theme:       getTheme(r),
		Extra:       string(extra),
	}

	err = tpl.ExecuteTemplate(w, "index"+templateSuffix, state)
	if err != nil {
		return err
	}

	return nil
}

func Render(templateName string, w http.ResponseWriter, r *http.Request) error {
	return RenderExtra(templateName, w, r, nil)
}

// Render403 renders the default 403 error page
func Render403(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(403)
	return RenderBasic("403", w, r)
}
