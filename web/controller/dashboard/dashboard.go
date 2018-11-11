package dashboard

import (
	"log"
	"net/http"

	"github.com/pajlada/pajbot2/web/router"
	"github.com/pajlada/pajbot2/web/views"
)

func Load() {
	router.Get("/dashboard", Dashboard)
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	err := views.Render("dashboard", w, r)
	if err != nil {
		log.Println("Error rendering dashboard view:", err)
	}
}
