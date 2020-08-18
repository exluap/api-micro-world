/**
 * Project api-microworld created by exluap
 * Date: 14.08.2020 12:46
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package user

import (
	"github.com/exluap/api-microworld/internal/database"
	"github.com/exluap/api-microworld/internal/utils"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
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
