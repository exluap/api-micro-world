/**
 * Project api-microworld created by exluap
 * Date: 28.08.2020 14:39
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package users

import (
	"encoding/json"
	"github.com/exluap/api-microworld/internal/database"
	"github.com/exluap/api-microworld/internal/user"
	"github.com/exluap/api-microworld/internal/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

type resultOfUsers struct {
	Result bool        `json:"result"`
	Users  []user.User `json:"users"`
} //@name UserList

// GetListOfUser godoc
// @Description Getting all users in system
// @Summary Getting all users in system
// @Produce json
// @Router /admin/users [get]
// @Success 200 {object} users.resultOfUsers "result"
func GetListOfUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	db := database.GetDb()
	// добавляем логику на будущую проверку по ролям
	//ctx := r.Context()

	var resultOfUsers resultOfUsers

	err := db.Table("users").Find(&resultOfUsers.Users).Error

	if err != nil {
		logrus.Errorf("can not get all users from database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "can not get all users",
		}
		message.Respond(w)
		return
	}
	resultOfUsers.Result = true
	result, err := json.Marshal(resultOfUsers)

	if err != nil {
		logrus.Errorf("can not prepare response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "can not prepare response",
		}
		message.Respond(w)
		return
	}

	_, err = w.Write(result)

	if err != nil {
		logrus.Errorf("can not prepare response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "can not prepare response",
		}
		message.Respond(w)
		return
	}
}
