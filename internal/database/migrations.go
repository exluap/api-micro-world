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
	"time"
)

type User struct {
	gorm.Model
	UUID     uuid.UUID `swaggerignore:"true"`
	Email    string    `json:"email" example:"test@test.com"`
	Login    string    `json:"login" example:"test"`
	Password string    `json:"password" example:"123456"`
	LastAuth time.Time `json:"last_auth" gorm:"timestamp with time zone"`
} //@name User
