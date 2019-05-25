package banphrases

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pajbot/pajbot2/pkg/web/router"
	"github.com/pajbot/pajbot2/pkg/web/views"
)

func Load() {
	router.Get("/banphrases", handleBanphrases)
}

func handleBanphrases(w http.ResponseWriter, r *http.Request) {
	banphrases := []string{
		"a", "b", "c",
	}
	extra, _ := json.Marshal(banphrases)

	err := views.RenderExtra("banphrases", w, r, extra)
	if err != nil {
		log.Println("Error rendering banphrases view:", err)
	}
}
