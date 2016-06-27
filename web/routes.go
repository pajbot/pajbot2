package web

import (
	"io/ioutil"
	"net/http"
)

func (boss *Boss) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./web/models/dashboard.html")
	if err != nil {
		w.Write([]byte("error: " + err.Error()))
	} else {
		w.Write(data)
	}
}
