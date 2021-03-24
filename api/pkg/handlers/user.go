package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/contractGuru/pkg/db"
	"github.com/contractGuru/pkg/logger"
	"github.com/contractGuru/pkg/models"
	"github.com/contractGuru/pkg/secure"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// LoginUser authenticates a user.
func LoginUser(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-type") == "application/json" {
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
		} else {
			respondError(w, http.StatusBadRequest, "Invalid content-type")
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

func GetUserInfo(db *db.DB) http.HandlerFunc {
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

		type userInfo struct {
			UserID     int `json:"user_id"`
			UserRoleID int `json:"user_role_id"`
			BuID       int `json:"bu_id"`
		}

		userinfo := &userInfo{}

		params := mux.Vars(r)["id"]
		userId, err := strconv.Atoi(params)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := db.Grm.Debug().
			Raw(`SELECT user.user_id, role.user_role_id, bu.bu_id
				FROM user, user_has_user_role AS role, user_has_bu AS bu
				WHERE user.user_id=role.user_id
				AND user.user_id=bu.user_id
				AND user.user_id=?`, userId).
			Scan(&userinfo).
			Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, userinfo)
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
