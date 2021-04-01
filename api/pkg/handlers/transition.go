package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/contractGuru/pkg/db"
	"github.com/contractGuru/pkg/models"
	"github.com/gorilla/mux"
)

func CreateRequestTransition(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-type") == "application/json" {
			session, err := db.Store.Get(r, "session")
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			if session.Values["auth"] != true {
				respondError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			newTransition := &models.RequestTransition{}

			if err = json.NewDecoder(r.Body).Decode(newTransition); err != nil {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}

			userId := session.Values["id"].(int)

			params := mux.Vars(r)["id"]
			requestId, err := strconv.Atoi(params)
			if err != nil {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}

			tx := db.Grm.Begin()

			if err := tx.Debug().
				Exec(`INSERT INTO request_transition (request_id, transition_id, user_id, created_at)
				VALUES(?, ?, ?, ?)`, requestId, newTransition.TransitionID, userId, time.Now()).
				Error; err != nil {
				tx.Rollback()
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			if err := tx.Debug().
				Table("request").
				Where("request_id=?", requestId).
				Update("current_state_id", newTransition.EndStateID).
				Error; err != nil {
				tx.Rollback()
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			tx.Commit()

			respondJSON(w, http.StatusCreated, "")
		} else {
			respondError(w, http.StatusBadRequest, "Invalid content-type")
			return
		}
	}
}
