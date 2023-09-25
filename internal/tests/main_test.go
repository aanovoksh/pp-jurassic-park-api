package tests

import (
	"log"
	"os"
	"pp-jurassic-park-api/internal/api/handlers"
	apimodels "pp-jurassic-park-api/internal/api/models"
	"pp-jurassic-park-api/internal/db"
	dbmodels "pp-jurassic-park-api/internal/db/models"
	"testing"

	"github.com/gin-gonic/gin"
)

var activeCage dbmodels.Cage
var downCage dbmodels.Cage
var cageToBeUpdated dbmodels.Cage
var cageToBeRemoved dbmodels.Cage

var cageWithTyrannosaurus dbmodels.Cage
var cageWithSpinosaurus dbmodels.Cage
var cageWithTriceratops dbmodels.Cage
var cageWithStegosaurus dbmodels.Cage

var tyrannosaurus dbmodels.Dinosaur
var spinosaurus dbmodels.Dinosaur
var triceratops dbmodels.Dinosaur
var stegosaurus dbmodels.Dinosaur
var dinosaurToBeRemoved dbmodels.Dinosaur

var cageIDsToCleanup []uint
var dinosaurIDsToCleanup []uint

var router *gin.Engine

func TestMain(m *testing.M) {
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations for test database: %v", err)
	}

	activeCage = CreateTestCage(2, apimodels.Active)
	downCage = CreateTestCage(2, apimodels.Down)
	cageToBeUpdated = CreateTestCage(3, apimodels.Active)
	cageToBeRemoved = CreateTestCage(4, apimodels.Down)

	cageWithTyrannosaurus = CreateTestCage(3, apimodels.Active)
	cageWithSpinosaurus = CreateTestCage(3, apimodels.Active)
	cageWithTriceratops = CreateTestCage(3, apimodels.Active)
	cageWithStegosaurus = CreateTestCage(3, apimodels.Active)

	tyrannosaurus = CreateTestDinosaur("Terry", apimodels.Tyrannosaurus, apimodels.Carnivore, cageWithTyrannosaurus.ID)
	spinosaurus = CreateTestDinosaur("Sally", apimodels.Spinosaurus, apimodels.Carnivore, cageWithSpinosaurus.ID)
	triceratops = CreateTestDinosaur("Tony", apimodels.Triceratops, apimodels.Herbivore, cageWithTriceratops.ID)
	stegosaurus = CreateTestDinosaur("Sony", apimodels.Stegosaurus, apimodels.Herbivore, cageWithStegosaurus.ID)
	dinosaurToBeRemoved = CreateTestDinosaur("Barry", apimodels.Brachiosaurus, apimodels.Herbivore, activeCage.ID)

	router = setupRouter()

	os.Exit(m.Run())

	DeleteTestDinosaurs(dinosaurIDsToCleanup)
	DeleteTestCages(cageIDsToCleanup)
}

func init() {
	os.Setenv("GO_ENV", "test")
	gin.SetMode(gin.TestMode)
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/cages", handlers.GetCages)
	router.GET("/cages/:id", handlers.GetCage)
	router.POST("/cages", handlers.CreateCage)
	router.PATCH("/cages/:id", handlers.UpdateCagePowerStatus)
	router.DELETE("/cages/:id", handlers.DeleteCage)

	router.GET("/dinosaurs", handlers.GetDinosaurs)
	router.GET("/dinosaurs/:id", handlers.GetDinosaur)
	router.POST("/dinosaurs", handlers.AddDinosaur)
	router.PATCH("/dinosaurs/:id", handlers.MoveDinosaur)
	router.DELETE("/dinosaurs/:id", handlers.RemoveDinosaur)
	return router
}
