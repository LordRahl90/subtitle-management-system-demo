package main

import (
	"fmt"
	"log"
	"os"

	"translations/servers"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}
	signingSecret := os.Getenv("SIGNING_SECRET")

	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}

	server, err := servers.New(db, signingSecret, "outputs/")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Router.Run(":8080"))
}

func setupDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
