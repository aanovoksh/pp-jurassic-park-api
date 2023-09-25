package main

import (
	"log"
	"pp-jurassic-park-api/internal/db"
	"pp-jurassic-park-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	_, err := db.Connect()
	if err != nil {
		log.Fatal("Could not connect to the database:", err)
	}

	engine := gin.Default()

	engine.GET("/tests", handlers.Test)

	err = engine.Run()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
