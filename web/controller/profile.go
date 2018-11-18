package controller

import (
	"log"
	"net/http"

	"github.com/pajlada/pajbot2/web/views"
)

func handleProfile(w http.ResponseWriter, r *http.Request) {
	err := views.Render("profile", w, r)
	if err != nil {
		log.Println("Error rendering dashboard view:", err)
	}
}
