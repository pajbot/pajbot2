package static

import (
	"net/http"
	"path/filepath"

	"github.com/pajbot/pajbot2/internal/config"
	"github.com/pajbot/pajbot2/pkg/web/router"
)

func Load() {
	staticPath := filepath.Join(config.WebStaticPath, "static")

	// Serve files statically from ./web/static in /static
	router.PathPrefix("/static", http.StripPrefix("/static", http.FileServer(http.Dir(staticPath))))
}
