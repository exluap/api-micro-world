/**
 * Project api-microworld created by exluap
 * Date: 02.08.2020 00:32
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package user

import (
	"encoding/json"
	"github.com/exluap/api-microworld/internal/database"
	"github.com/exluap/api-microworld/internal/utils"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func validate(account database.User) (utils.Message, bool) {
	if !strings.Contains(account.Email, "@") {
		log.Warnf("email is empty or incorrect %s", account.Email)
		return utils.Message{Result: false, Message: "Email is empty or incorrect"}, false
	}

	if len(account.Password) < 6 {
		log.Warnf("Password is empty or <6 for user %s", account.Email+" "+account.Login)
		return utils.Message{Result: false, Message: "Password is required. Minimum 6 elements "}, false
	}

	//Email должен быть уникальным
	temp := &database.User{}

	//проверка на наличие ошибок и дубликатов электронных писем
	err := database.GetDb().Table("users").Where("email = ? or login = ?", account.Email, account.Login).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("can not connect to database or find user %v, with error: %v", account, err)
		return utils.Message{Result: false, Message: "Connection error. Please retry"}, false
	}
	if temp.Email != "" {
		log.Warn("user or email existing")
		return utils.Message{Result: false, Message: "Email address  or login already in use by another user."}, false
	}

	return utils.Message{Result: false, Message: "Requirement passed"}, true
}

// RegisterUser godoc
// @Description Register user with data. ATTENTION! Password must be >=6 symbols
// @Summary Register user with specified model
// @Produce json
// @Router /user/register [post]
// @Success 200 {object} utils.Message "result and token"
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Errorf("error with read body at register query: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "error with read body",
		}
		message.Respond(w)
		return
	}

	var user database.User
	err = json.Unmarshal(body, &user)

	if err != nil {
		log.Errorf("can not unmarshall body in register query: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "error with unmarshal body",
		}
		message.Respond(w)
		return
	}

	db := database.GetDb()

	user.UUID, err = uuid.NewUUID()

	if err != nil {
		log.Errorf("error with create uuid for user %s: %v", user.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "can not create user",
		}
		message.Respond(w)
		return
	}

	if res, ok := validate(user); !ok {
		w.WriteHeader(http.StatusConflict)
		res.Respond(w)
		return
	}

	err = db.Create(&user).Error

	if err != nil {
		log.Errorf("can not create user %s at database: %v", user.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "can not create user",
		}
		message.Respond(w)
		return
	}

	log.Infof("user %s successful register, his uuid %s!", user.Email, user.UUID)
	w.WriteHeader(http.StatusOK)
	message := utils.Message{
		Result:  true,
		Message: utils.GenerateToken(user.UUID),
	}
	message.Respond(w)
	return
}

// AuthUser godoc
// @Summary Authentication user
// @Produce json
// @Router /user/login [post]
// @Success 200 {object} utils.Message "result and token"
func AuthUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Errorf("can not read body in login query for user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "error with read body",
		}
		message.Respond(w)
		return
	}

	var temp database.User

	err = json.Unmarshal(body, &temp)

	if err != nil {
		log.Errorf("can not unmarshall body in login query for user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "can not unmarshall body",
		}
		message.Respond(w)
		return
	}

	user := &database.User{}

	err = database.GetDb().Table("users").Where("(login = ? or email = ? ) and password = ?",
		temp.Login, temp.Email, temp.Password).First(user).Error

	if err != nil {
		log.Errorf("can not find user %s or password is incorrect:  %v", temp.Login+" "+temp.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "can not find user or password is incorrect",
		}
		message.Respond(w)
		return
	}

	user.LastAuth = time.Now()

	err = database.GetDb().Table("users").Where("uuid = ?", user.UUID).Save(user).Error

	if err != nil {
		log.Errorf("can not save user %v last login with error: %v", user.UUID, err)
		w.WriteHeader(http.StatusInternalServerError)
		message := utils.Message{
			Result:  false,
			Message: "can not save user auth",
		}
		message.Respond(w)
		return
	}

	log.Infof("user %s successful login", user.UUID)
	w.WriteHeader(http.StatusOK)
	message := utils.Message{
		Result:  true,
		Message: utils.GenerateToken(user.UUID),
	}
	message.Respond(w)
	return
}
