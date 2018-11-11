package home

import (
	"log"
	"net/http"

	"github.com/pajlada/pajbot2/web/router"
	"github.com/pajlada/pajbot2/web/views"
)

func Load() {
	router.Get("/", Home)
	router.Get("/home", Home)
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := views.Render("home", w, r)
	if err != nil {
		log.Println("Error rendering dashboard view:", err)
	}
}
