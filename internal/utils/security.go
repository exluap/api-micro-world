/**
 * Project api-microworld created by exluap
 * Date: 14.08.2020 19:15
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package utils

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

type Token struct {
	UUID uuid.UUID
	jwt.StandardClaims
}

func GenerateToken(uuid uuid.UUID) string {
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

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization") //Получение токена

		if tokenHeader == "" { //Токен отсутствует, возвращаем  403 http-код Unauthorized
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			message := Message{
				Result:  false,
				Message: "Missing auth token",
			}
			message.Respond(w)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //Токен обычно поставляется в формате `Bearer {token-body}`, мы проверяем, соответствует ли полученный токен этому требованию
		if len(splitted) != 2 {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			message := Message{
				Result:  false,
				Message: "Invalid/Malformed auth token",
			}
			message.Respond(w)
			return
		}

		tokenPart := splitted[1] //Получаем вторую часть токена
		tk := &Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_PASSWORD")), nil
		})

		if err != nil { //Неправильный токен, как правило, возвращает 403 http-код
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			message := Message{
				Result:  false,
				Message: "Malformed authentication token",
			}
			message.Respond(w)
			return
		}

		if !token.Valid { //токен недействителен, возможно, не подписан на этом сервере
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			message := Message{
				Result:  false,
				Message: "Token is not valid.",
			}
			message.Respond(w)
			return
		}

		//Всё прошло хорошо, продолжаем выполнение запроса
		ctx := context.WithValue(r.Context(), "userID", tk.UUID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //передать управление следующему обработчику!
	})
}
