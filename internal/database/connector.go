/**
 * Project api-microworld created by exluap
 * Date: 02.08.2020 15:17
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package database

import (
	"github.com/jinzhu/gorm"
	"log"
	"os"
)
import _ "github.com/jinzhu/gorm/dialects/postgres"

var Db *gorm.DB

func init() {
	db, err := gorm.Open("postgres", "host="+os.Getenv("DB_HOST")+" port="+os.Getenv("DB_PORT")+
		" user="+os.Getenv("DB_USER")+" dbname="+os.Getenv("DB_NAME")+" password="+os.Getenv("DB_PASSWORD")+
		" sslmode=disable")

	if err != nil {
		log.Fatal("can not open database connect: ", err)
	}

	db.AutoMigrate(&User{})

	Db = db

}

func GetDb() *gorm.DB {
	return Db
}
