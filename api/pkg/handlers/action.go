package handlers

import (
	"net/http"
	"strconv"

	"github.com/contractGuru/pkg/db"
	"github.com/contractGuru/pkg/models"
	"github.com/gorilla/mux"
)

// GetPendingRequestActions queries for the next action to be taken.
func GetPendingRequestActions(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := db.Store.Get(r, "session")
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if session.Values["auth"] != true {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		action := &[]models.Action{}

		bu := session.Values["bu"].(int)
		role := session.Values["role"].(int)

		params := mux.Vars(r)["id"]
		requestId, err := strconv.Atoi(params)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := db.Grm.Debug().
			Raw(`SELECT ta.action_id, request.finance_flag
				FROM request, transition, transition_action AS ta
				WHERE request.current_state_id=transition.start_state_id
				AND request.bu_id=?
				AND transition.user_role_id=?
				AND request.request_id=?
				AND transition.transition_id=ta.transition_id`, bu, role, requestId).
			Scan(&action).
			Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(*action) == 0 {
			respondError(w, http.StatusNoContent, "")
			return
		}

		respondJSON(w, http.StatusOK, action)
	}
}
