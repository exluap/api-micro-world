/**
 * Project api-microworld created by exluap
 * Date: 02.08.2020 18:54
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

type Token struct {
	UUID uuid.UUID
	jwt.StandardClaims
}

type User struct {
	Login    string    `json:"login"`
	Email    string    `json:"email"`
	LastAuth time.Time `json:"last_auth" `
} //@name User
