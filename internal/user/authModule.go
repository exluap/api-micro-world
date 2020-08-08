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
	"github.com/dgrijalva/jwt-go"
	"github.com/exluap/api-microworld/internal/database"
	"github.com/exluap/api-microworld/internal/utils"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func validate(account database.User) (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		log.Warnf("email is empty or incorrect %s", account.Email)
		return utils.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		log.Warnf("Password is empty or <6 for user %s", account.Email+" "+account.Login)
		return utils.Message(false, "Password is required. Minimum 6 elements "), false
	}

	//Email должен быть уникальным
	temp := &database.User{}

	//проверка на наличие ошибок и дубликатов электронных писем
	err := database.GetDb().Table("users").Where("email = ? or login = ?", account.Email, account.Login).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("can not connect to database or find user %v, with error: %v", account, err)
		return utils.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		log.Warn("user or email existing")
		return utils.Message(false, "Email address  or login already in use by another user."), false
	}

	return utils.Message(false, "Requirement passed"), true
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Errorf("error with read body at register query: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.Message(false, "error with read body"))
		return
	}

	var user database.User
	err = json.Unmarshal(body, &user)

	if err != nil {
		log.Errorf("can not unmarshall body in register query: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.Message(false, "error with unmarshall body"))
		return
	}

	db := database.GetDb()

	user.UUID, err = uuid.NewUUID()

	if err != nil {
		log.Errorf("error with create uuid for user %s: %v", user.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.Message(false, "error with creating user (generate uuid)"))
		return
	}

	if res, ok := validate(user); !ok {
		w.WriteHeader(http.StatusConflict)
		utils.Respond(w, res)
		return
	}

	err = db.Create(&user).Error

	if err != nil {
		log.Errorf("can not create user %s at database: %v", user.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.Message(false, "error with creating user at database"))
		return
	}

	log.Infof("user %s successful register, his uuid %s!", user.Email, user.UUID)
	w.WriteHeader(http.StatusOK)
	utils.Respond(w, utils.Message(true, generateToken(user.UUID)))
	return
}

func AuthUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Errorf("can not read body in login query for user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.Message(false, "error with read body"))
		return
	}

	var temp database.User

	err = json.Unmarshal(body, &temp)

	if err != nil {
		log.Errorf("can not unmarshall body in login query for user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.Message(false, "can not unmarshall body"))
		return
	}

	user := &database.User{}

	err = database.GetDb().Table("users").Where("(login = ? or email = ? ) and password = ?",
		temp.Login, temp.Email, temp.Password).First(user).Error

	if err != nil {
		log.Errorf("can not find user %s or password is incorrect:  %v", temp.Login+" "+temp.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.Message(false, "can not find user or password is incorrect"))
		return
	}

	log.Infof("user %s successful login", user.UUID)
	w.WriteHeader(http.StatusOK)
	utils.Respond(w, utils.Message(true, generateToken(user.UUID)))
	return
}

func generateToken(uuid uuid.UUID) string {
	tokenClaims := &Token{
		UUID: uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 3).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))

	if err != nil {
		log.Errorf("error with generating jwt key for user: %s", uuid)
	}

	return tokenString
}
