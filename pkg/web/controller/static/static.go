package static

import (
	"net/http"

	"github.com/pajbot/pajbot2/pkg/web/router"
)

func Load() {
	// Serve files statically from ./web/static in /static
	router.PathPrefix("/static", http.StripPrefix("/static", http.FileServer(http.Dir("../../web/static/"))))
}
