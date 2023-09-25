package main

import (
	"log"
	"pp-jurassic-park-api/internal/db"
	"pp-jurassic-park-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Run DB Migration
	err := db.Migrate()
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	router := gin.Default()

	// Cages API
	router.GET("/cages", handlers.GetCages)
	router.GET("/cages/:id", handlers.GetCage)
	router.POST("/cages", handlers.CreateCage)
	router.PATCH("/cages/:id", handlers.UpdateCagePowerStatus)
	router.DELETE("/cages/:id", handlers.DeleteCage)

	// Dinosaur API
	router.GET("/dinosaurs", handlers.GetDinosaurs)
	router.GET("/dinosaurs/:id", handlers.GetDinosaur)
	router.POST("/dinosaurs", handlers.AddDinosaur)
	router.PATCH("/dinosaurs/:id", handlers.MoveDinosaur)
	router.DELETE("/dinosaurs/:id", handlers.RemoveDinosaur)

	// Start server
	err = router.Run()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
