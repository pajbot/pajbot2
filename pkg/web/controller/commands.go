package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pajbot/pajbot2/pkg/commandlist"
	"github.com/pajbot/pajbot2/pkg/web/views"
)

func handleCommands(w http.ResponseWriter, r *http.Request) {
	commands := commandlist.Commands()
	bytes, err := json.Marshal(commands)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = views.RenderExtra("commands", w, r, bytes)
	if err != nil {
		log.Println("Error rendering dashboard view:", err)
	}
}
