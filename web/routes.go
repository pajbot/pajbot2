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

type dashboardData struct {
	WSHost string
}

func (b *Boss) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	// tpl.Execute(w, "dashboard.html", dashboardData{
	// 	WSHost: b.WSHost,
	// })

	// the template thing didnt work because ng also uses {{ }}
	// and im too tired to fix it LUL
	http.ServeFile(w, r, "./web/models/dashboard.html")
}
