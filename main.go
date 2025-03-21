package main

import (
	"kreditplus/config"
	"kreditplus/internal/domain"
	"kreditplus/internal/route"
	"log"
)

func main() {
	config.ConnectDB()
	config.DB.AutoMigrate(&domain.User{}, &domain.Customer{}, &domain.Limit{}, &domain.Transaction{})

	var count int64
	config.DB.Model(&domain.User{}).Where("user_username = ?", "admin").Count(&count)
	if count == 0 {
		password := "$2a$12$kcu7bfXuDBZn0r9hN50xH.qQJynE9f6gUAq9Jx0diDeQA97AfT9P2"
		admin := domain.User{
			UserUsername: "admin",
			UserPassword: password,
			UserRole:     "admin",
		}
		config.DB.Create(&admin)
		log.Println("Default admin user created")
	} else {
		log.Println("Admin user already exists")
	}

	r := route.SetupRouter()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	log.Println("Server running on port 8080")
	r.Run(":8080")
}
