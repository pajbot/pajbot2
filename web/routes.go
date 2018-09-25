package web

import (
	"fmt"
	"net/http"

	"github.com/gernest/hot"
)

var tpl *hot.Template

func init() {
	var err error
	tpl, err = hot.New(&hot.Config{
		Watch:          true,
		BaseName:       "base",
		Dir:            "web/models/",
		FilesExtension: []string{".html"},
		LeftDelim:      "[[",
		RightDelim:     "]]",
	})
	if err != nil {
		panic(err)
	}
}

type dashboardData struct {
	WSHost string
}

func (b *Boss) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	fmt.Printf("xd %#v\n", b.WSHost)
	tpl.Execute(w, "dashboard.html", dashboardData{
		WSHost: b.WSHost,
	})
}
