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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func validate(account database.User) (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return utils.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return utils.Message(false, "Password is required"), false
	}

	//Email должен быть уникальным
	temp := &database.User{}

	//проверка на наличие ошибок и дубликатов электронных писем
	err := database.GetDb().Table("users").Where("email = ? or login = ?", account.Email, account.Login).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return utils.Message(false, "Email address  or login already in use by another user."), false
	}

	return utils.Message(false, "Requirement passed"), true
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("error with read body at register func: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.Message(false, "error with read body"))
		return
	}

	var user database.User
	err = json.Unmarshal(body, &user)

	if err != nil {
		log.Printf("error with unmarshall body reg: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.Message(false, "error with unmarshall body"))
		return
	}

	db := database.GetDb()

	user.UUID, err = uuid.NewUUID()

	if err != nil {
		log.Printf("error with create uuid for user: %v", err)
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
		log.Printf("can not create user at database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.Message(false, "error with creating user at database"))
		return
	}

	tk := &Token{UUID: user.UUID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))

	w.WriteHeader(http.StatusOK)
	utils.Respond(w, utils.Message(true, tokenString))
	return
}

func AuthUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("can not read body:  %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.Message(false, "error with read body"))
		return
	}

	var user database.User

	err = json.Unmarshal(body, &user)

	if err != nil {
		log.Printf("can not unmarshall body:  %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.Message(false, "can not unmarshall body"))
		return
	}

	temp := &database.User{}

	err = database.GetDb().Table("users").Where("(login = ? or email = ? ) and password = ?",
		user.Login, user.Email, user.Password).First(temp).Error

	if err != nil {
		log.Printf("can not find user:  %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.Message(false, "can not find user or password is incorrect"))
		return
	}

	tk := &Token{UUID: temp.UUID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))

	w.WriteHeader(http.StatusOK)
	utils.Respond(w, utils.Message(true, tokenString))
	return
}
