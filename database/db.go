package database

import (
	"github.com/salamalfis/Assigment2golang/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Alfis"
	dbname   = "order_by"
)

var (
	db  *gorm.DB
	err error
)

func StartDB() {
	config := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		log.Fatal("errors")

	}

	fmt.Println("successful")
	db.Debug().AutoMigrate(models.Orders{}, models.Items{})

}
func GetDB() *gorm.DB {
	return db
}
