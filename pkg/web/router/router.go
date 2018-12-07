package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// r will most likely need to be mutex locked

var (
	r *mux.Router
)

func init() {
	r = mux.NewRouter()
}

func Subrouter(path string) *mux.Router {
	return r.PathPrefix(path).Subrouter()
}

func RGet(R *mux.Router, path string, handler http.HandlerFunc) *mux.Route {
	return R.HandleFunc(path, handler).Methods("GET")
}

func RPost(R *mux.Router, path string, handler http.HandlerFunc) *mux.Route {
	return R.HandleFunc(path, handler).Methods("POST")
}

func Get(path string, handler http.HandlerFunc) {
	RGet(r, path, handler)
}

func PathPrefix(path string, handler http.Handler) {
	r.PathPrefix(path).Handler(handler)
}

func RHandleFunc(R *mux.Router, path string, handler http.HandlerFunc) {
	R.HandleFunc(path, handler)
}

func HandleFunc(path string, handler http.HandlerFunc) {
	RHandleFunc(r, path, handler)
}

func Instance() *mux.Router {
	return r
}

func Debug() {
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		fmt.Println(t)
		return nil
	})
}
