package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/contractGuru/pkg/db"
)

//CreateCompletedContractData handles inserting new completed contract entry
func CreateCompletedContractData(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-type") == "application/json" {
			session, err := db.Store.Get(r, "session")
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			if session.Values["auth"] != true {
				respondError(w, http.StatusUnauthorized, err.Error())
				return
			}

			newCompletedData = &models.Completed{}

			if err = json.NewDecoder(r.Body).Decode(newCompletedData); err != nil {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}

			if err := db.Grm.Debug().
				Table("completed").
				Create(&newCompletedData).
				Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondJson(w, http.StatusOK, nil)
		} else {
			respondError(w, http.StatusBadRequest, "Invalid content type")
			return
		}
	}
}
