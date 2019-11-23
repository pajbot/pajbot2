package channel

import (
	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg/web/controller/api/channel/banphrases"
	"github.com/pajbot/pajbot2/pkg/web/controller/api/channel/moderation"
)

// Load routes for /api/channel/:channel_id
func Load(parent *mux.Router) {
	m := parent.PathPrefix(`/channel/{channel_id:\w+}`).Subrouter()

	// Load subroutes for /api/channel/:channel_id/moderation/
	moderation.Load(m)

	// Load subroutes for /api/channel/:channel_id/banphrases/
	banphrases.Load(m)
}
