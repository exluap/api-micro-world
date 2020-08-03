/**
 * Project api-microworld created by exluap
 * Date: 02.08.2020 00:39
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

// Миграции в БД для дальнейшей работы

package database

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UUID     uuid.UUID
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
