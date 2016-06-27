package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pajlada/pajbot2/common"
)

// Boss xD
type Boss struct {
	Host string
}

// Init returns a webBoss which hosts the website
func Init(config *common.Config) *Boss {
	webHost := ":2355" // TEMPORARY
	boss := &Boss{
		Host: webHost,
	}
	return boss
}

// Run xD
func (boss *Boss) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", boss.wsHandler)
	r.HandleFunc("/", boss.rootHandler)
	r.HandleFunc("/dashboard", boss.dashboardHandler)

	log.Infof("Starting web on host %s", boss.Host)
	http.ListenAndServe(boss.Host, r)
}

func (boss *Boss) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "<h1>xD</h1>")
}

func (boss *Boss) wsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: readd this once we have a proper config shit
	/*
		if r.Header.Get("Origin") != "http://"+r.Host {
			http.Error(w, "Origin not allowed", 403)
			return
		}
	*/
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	go wsHandler(conn)
}
