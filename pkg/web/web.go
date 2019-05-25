package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/pajbot/pajbot2/pkg/common/config"
	"github.com/pajbot/pajbot2/pkg/web/router"
)

func Run(cfg *config.WebConfig) {
	fmt.Printf("Starting web on host %s\n", cfg.Host)
	corsObj := handlers.AllowedOrigins([]string{"*"})
	err := http.ListenAndServe(cfg.Host, handlers.CORS(corsObj)(router.Instance()))
	if err != nil {
		log.Fatal(err)
	}
}
