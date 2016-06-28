package web

import (
	"net/http"

	"github.com/gernest/hot"
)

var tpl, _ = hot.New(&hot.Config{
	Watch:          true,
	BaseName:       "base",
	Dir:            "web/models/",
	FilesExtension: []string{".html"},
})

func (boss *Boss) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	tpl.Execute(w, "dashboard.html", r.Host)
}
