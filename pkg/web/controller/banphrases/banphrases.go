package banphrases

import (
	"log"
	"net/http"

	"github.com/pajlada/pajbot2/pkg/web/router"
	"github.com/pajlada/pajbot2/pkg/web/views"
)

func Load() {
	router.Get("/banphrases", handleBanphrases)
}

func handleBanphrases(w http.ResponseWriter, r *http.Request) {
	banphrases := []string{
		"a", "b", "c",
	}
	err := views.RenderExtra("banphrases", w, r, banphrases)
	if err != nil {
		log.Println("Error rendering banphrases view:", err)
	}
}
