package handlers

import (
	"net/http"

	"github.com/contractGuru/pkg/db"
	"github.com/contractGuru/pkg/models"
)

func GetPendingRequests(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := db.Store.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if session.Values["auth"] != true {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		request := []models.Request{}

		db.Grm.Debug().Table("request").Find(&request)

		respondJSON(w, http.StatusOK, request)
	}
}