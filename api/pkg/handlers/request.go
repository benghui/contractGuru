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
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		request := []models.Request{}

		if err := db.Grm.Debug().
			Table("request").
			Find(&request).
			Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, request)
	}
}

// CreateRequest handles inserting new contract request entry.
func CreateRequest(db *db.DB) http.HandlerFunc {
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
	}
}
