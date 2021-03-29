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

			userId := userData["id"].(int)

			userinfo, err := getUserRoleAndBu(db, userId)

			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			session.Values["id"] = userData["id"]
			session.Values["auth"] = true
			session.Values["role"] = userinfo.UserRoleID
			session.Values["bu"] = userinfo.BuID

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

// getUserRoleAndBu queries for user role and bu.
func getUserRoleAndBu(db *db.DB, userId int) (*models.UserInfo, error) {

	userinfo := &models.UserInfo{}

	if err := db.Grm.Debug().
		Raw(`SELECT user.user_id, role.user_role_id, bu.bu_id
			FROM user, user_has_user_role AS role, user_has_bu AS bu
			WHERE user.user_id=role.user_id
			AND user.user_id=bu.user_id
			AND user.user_id=?`, userId).
		Scan(&userinfo).
		Error; err != nil {
		logger.Error.Println(err.Error())
		return nil, err
	}
	return userinfo, nil
}
