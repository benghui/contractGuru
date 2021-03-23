package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/contractGuru/pkg/db"
	"github.com/contractGuru/pkg/logger"
	"github.com/contractGuru/pkg/models"
	"github.com/contractGuru/pkg/secure"
	"golang.org/x/crypto/bcrypt"
)

// LoginUser authenticates a user.
func LoginUser(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := db.Store.Get(r, "session")
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		user := models.User{}

		if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		user.Prepare()

		userData, err := passwordCheck(db, user.Username, user.Password)

		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		session.Values["id"] = userData["id"]
		session.Values["auth"] = true

		if err = session.Save(r, w); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

// LogoutUser revokes session.
func LogoutUser(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := db.Store.Get(r, "session")
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		session.Values["id"] = nil
		session.Values["auth"] = nil
		session.Options.MaxAge = -1

		if err = session.Save(r, w); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

// passwordCheck verifies password.
func passwordCheck(db *db.DB, username, password string) (map[string]interface{}, error) {
	userData := make(map[string]interface{})

	user := models.User{}

	if err := db.Grm.Debug().
		Table("user").
		Model(models.User{}).
		Select("user_id", "username", "password").
		Where("username = ?", username).
		Take(&user).
		Error; err != nil {
		logger.Error.Println(err.Error())
		return nil, err
	}

	err := secure.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		logger.Error.Println(err.Error())
		return nil, err
	}

	userData["id"] = user.UserID
	userData["username"] = user.Username

	return userData, nil
}
