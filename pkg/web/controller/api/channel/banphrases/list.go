package banphrases

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/utils"
	"github.com/pajlada/pajbot2/pkg/web/state"
	"github.com/pajlada/pajbot2/pkg/webutils"
)

type banphrase struct {
	ID          string
	Enabled     bool
	Description string
	Phrase      string
}

type listResponse struct {
	Banphrases []banphrase
	ChannelID  string
}

func handleList(w http.ResponseWriter, r *http.Request) {
	c := state.Context(w, r)

	if !webutils.RequirePermission(w, c, pkg.PermissionModeration) {
		return
	}

	vars := mux.Vars(r)
	var response listResponse

	response.ChannelID = vars["channelID"]

	const queryF = "SELECT `id`, `enabled`, `description`, `phrase` FROM `Banphrase`"

	rows, err := c.SQL.Query(queryF)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var bp banphrase
		if err := rows.Scan(&bp.ID, &bp.Enabled, &bp.Description, &bp.Phrase); err != nil {
			panic(err)
		}

		response.Banphrases = append(response.Banphrases, bp)
	}

	utils.WebWrite(w, response)
}
