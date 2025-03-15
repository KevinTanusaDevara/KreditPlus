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

	r := route.SetupRouter()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	log.Println("Server running on port 8080")
	r.Run(":8080")
}
