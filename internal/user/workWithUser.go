/**
 * Project api-microworld created by exluap
 * Date: 14.08.2020 12:46
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package user

import (
	"encoding/json"
	"github.com/exluap/api-microworld/internal/database"
	"github.com/exluap/api-microworld/internal/utils"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// GetUserInfo godoc
// @Description Getting user's profile info
// @Summary Getting user's profile
// @Produce json
// @Param userUUID path string true "user uuid from token"
// @Router /user/{userUUID}/info [get]
// @Success 200 {object} utils.Message "user info"
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	userUUID := chi.URLParam(r, "userId")
	db := database.GetDb()
	ctx := r.Context()
	logrus.Infof("user %s request user's %s info", ctx.Value("userID"), userUUID)

	var user User

	err := db.Where("uuid = ?", userUUID).First(&user).Error

	if err != nil {
		logrus.Errorf("error with getting user %s info: %v", userUUID, err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "error with getting info of user",
		}
		message.Respond(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	message := utils.Message{
		Result:  true,
		Message: user,
	}
	message.Respond(w)
}

// DeleteUser godoc
// @summary Deleting user profile
// @produce json
// @param userUUID path string true "user uuid from token"
// @router /user/{userUUID} [delete]
// @sucess 200 {object} utils.Message "result of operation"
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	userUUID := chi.URLParam(r, "userId")
	db := database.GetDb()
	ctx := r.Context()
	logrus.Infof("user %s wanna delete user's %s profile", ctx.Value("userID"), userUUID)

	if userUUID == ctx.Value("userId") {
		w.WriteHeader(http.StatusBadRequest)
		message := utils.Message{
			Result:  false,
			Message: "you can not delete yourself",
		}
		message.Respond(w)
		return
	}

	var user database.User

	err := db.Delete(&user, "uuid = ?", userUUID).Error

	if err != nil {
		logrus.Errorf("error with deleting user %s profile: %v", userUUID, err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "error with deleting user",
		}
		message.Respond(w)
		return
	}

	message := utils.Message{
		Result:  true,
		Message: "user's " + userUUID + " profile deleted",
	}
	message.Respond(w)
}

//UpdateUserInfo godoc
// @summary Updating user info
// @produce json
// @router /user/me [post]
// @sucess 200 {object} utils.Message "result of operation"
func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	db := database.GetDb()
	ctx := r.Context()

	logrus.Infof("Ctx: %v", ctx)

	logrus.Infof("user %v start to update self profile", ctx.Value("userID"))

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logrus.Errorf("error with body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "error with body",
		}
		message.Respond(w)
		return
	}

	userUpdate := struct {
		User
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}{}

	err = json.Unmarshal(body, &userUpdate)

	if err != nil {
		logrus.Errorf("error unmarshall body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "error with unmarshal body",
		}
		message.Respond(w)
		return
	}

	var tempUser database.User

	err = db.Where("uuid = ?", ctx.Value("userID")).First(&tempUser).Error

	if err != nil {
		logrus.Errorf("error with getting user by uuid %v: %v", ctx.Value("userID"), err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "error with update user profile",
		}
		message.Respond(w)
		return
	}

	if userUpdate.NewPassword != "" {
		if userUpdate.OldPassword != tempUser.Password {
			logrus.Warnf("user %v try to update password! But it is not correct", ctx.Value("userID"))
			w.WriteHeader(http.StatusBadRequest)
			message := utils.Message{
				Result:  false,
				Message: "can not update user profile. Incorrect password",
			}
			message.Respond(w)
			return
		}

		tempUser.Password = userUpdate.NewPassword
	}

	if userUpdate.Email != "" {
		tempUser.Email = userUpdate.Email
	}

	if userUpdate.Login != "" {
		tempUser.Login = userUpdate.Login
	}

	err = db.Where("uuid = ?", ctx.Value("userID")).Model(&database.User{}).Updates(tempUser).Error

	if err != nil {
		logrus.Errorf("error with saving updates user by uuid %v: %v", ctx.Value("userID"), err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "error with update user profile",
		}
		message.Respond(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	message := utils.Message{
		Result:  true,
		Message: "all saved",
	}
	message.Respond(w)
}
