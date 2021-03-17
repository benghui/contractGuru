package handlers

import (
	"encoding/json"
	"io/ioutil"
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
		defer r.Body.Close()

		session, err := db.Store.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		user := models.Users{}

		if err = json.Unmarshal(body, &user); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		user.Prepare()

		userData, err := passwordCheck(db, user.Username, user.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		session.Values["id"] = userData["id"]
		session.Values["auth"] = true

		if err = session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func LogoutUser(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		session, err := db.Store.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["id"] = nil
		session.Values["auth"] = nil
		session.Options.MaxAge = -1

		if err = session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func passwordCheck(db *db.DB, username, password string) (map[string]interface{}, error) {

	userData := make(map[string]interface{})

	user := models.Users{}

	if err := db.Grm.Debug().
		Model(models.Users{}).
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
