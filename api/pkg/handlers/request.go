package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/contractGuru/pkg/db"
	"github.com/contractGuru/pkg/models"
)

// GetPendingRequests handles query for contract requests that are not yet completed.
func GetPendingRequests(db *db.DB) http.HandlerFunc {
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

		request := &[]models.Request{}
		bu := session.Values["bu"].(int)

		if err := db.Grm.Debug().
			Table("request").
			Where("current_state_id != ? AND bu_id=?", 4, bu).
			Limit(10).
			Find(&request).
			Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(*request) == 0 {
			respondError(w, http.StatusNoContent, "")
			return
		}

		respondJSON(w, http.StatusOK, request)
	}
}

// CreateRequest handles inserting new contract request entry.
func CreateRequest(db *db.DB) http.HandlerFunc {
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

			newRequest := &models.Request{}

			if err = json.NewDecoder(r.Body).Decode(newRequest); err != nil {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}

			newRequest.Prepare()

			if err := db.Grm.Debug().
				Table("request").
				Omit("RequestID").
				Create(&newRequest).
				Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondJSON(w, http.StatusOK, nil)
		} else {
			respondError(w, http.StatusBadRequest, "Invalid content-type")
			return
		}
	}
}
