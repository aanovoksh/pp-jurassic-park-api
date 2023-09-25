package handlers

import (
	"net/http"
	"strconv"
	"strings"

	apimodels "pp-jurassic-park-api/internal/api/models"
	transform "pp-jurassic-park-api/internal/api/transform"
	db "pp-jurassic-park-api/internal/db"
	dbmodels "pp-jurassic-park-api/internal/db/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetDinosaur returns single dinosaur for the requested id.
// Used to retrieve data around single dinosaur at the Jurassic Park.
func GetDinosaur(c *gin.Context) {
	idParam := c.Param("id")
	dinosaurID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid dinosaur ID."})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to connect to database."})
		return
	}

	var dinosaur dbmodels.Dinosaur
	if err := dbConn.First(&dinosaur, dinosaurID).Error; err != nil {
		c.JSON(http.StatusNotFound, apimodels.ErrorResponse{Error: "Dinosaur not found."})
		return
	}

	c.JSON(http.StatusOK, apimodels.GetDinosaurResponse{Dinosaur: transform.DinosaurToApi(dinosaur)})
}

// GetDinosaurs returns all dinosaurs matching provided filters.
// Used to retrieve data around all current dinosaur at the Jurassic Park.
func GetDinosaurs(c *gin.Context) {
	var req apimodels.GetDinosaursRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid input."})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to connect to database."})
		return
	}

	var dinosaurs []dbmodels.Dinosaur
	if len(req.FilteredSpecies) > 0 {
		if err := dbConn.Where("species IN ?", req.FilteredSpecies).Find(&dinosaurs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to retrieve dinosaurs with filtered species."})
			return
		}
	} else {
		if err := dbConn.Find(&dinosaurs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to retrieve dinosaurs."})
			return
		}
	}

	c.JSON(http.StatusOK, apimodels.GetDinosaursResponse{Dinosaurs: transform.DinosaursToApi(dinosaurs)})
}

// AddDinosaur adds a new dinosaurs to the cage.
// Used when dinosaur is imported to the Jurassic Park.
func AddDinosaur(c *gin.Context) {
	var req apimodels.AddDinosaurRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid input."})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to connect to database."})
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid name. Name cannot be blank."})
		return
	}

	knownSpecies, species, dinosaurType := apimodels.LookupSpeciesType(req.Species)
	if !knownSpecies {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Unknown species."})
		return
	}

	dinosaur := dbmodels.Dinosaur{
		Name:    req.Name,
		Species: string(species),
		Type:    string(dinosaurType),
		CageID:  req.CageID,
	}

	if !canBeMovedToCage(c, dbConn, dinosaur, req.CageID) {
		return
	}

	if err := dbConn.Create(&dinosaur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to add dinosaur."})
		return
	}

	c.JSON(http.StatusOK, apimodels.AddDinosaurResponse{Dinosaur: transform.DinosaurToApi(dinosaur)})
}

// MoveDinosaur moves existing dinosaurs to a different cage.
// Used to move dinosaurs around the Jurassic Park.
func MoveDinosaur(c *gin.Context) {
	idParam := c.Param("id")
	dinosaurID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid dinosaur ID."})
		return
	}

	var req apimodels.MoveDinosaurRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid input."})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to connect to database."})
		return
	}

	var dinosaur dbmodels.Dinosaur
	if err := dbConn.First(&dinosaur, dinosaurID).Error; err != nil {
		c.JSON(http.StatusNotFound, apimodels.ErrorResponse{Error: "Dinosaur not found."})
		return
	}

	if dinosaur.CageID != req.CageID {
		if !canBeMovedToCage(c, dbConn, dinosaur, req.CageID) {
			return
		}
		dinosaur.CageID = req.CageID
		if err := dbConn.Save(&dinosaur).Error; err != nil {
			c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to move dinosaur."})
			return
		}
	}

	c.JSON(http.StatusOK, apimodels.MoveDinosaurResponse{Dinosaur: transform.DinosaurToApi(dinosaur)})
}

// RemoveDinosaur removes dinosaur from their existing cage.
// Used when dinosaur is exported from the Jurassic Park.
func RemoveDinosaur(c *gin.Context) {
	idParam := c.Param("id")
	dinosaurID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid dinosaur ID."})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to connect to database."})
		return
	}

	var dinosaur dbmodels.Dinosaur
	if err := dbConn.First(&dinosaur, dinosaurID).Error; err != nil {
		c.JSON(http.StatusNotFound, apimodels.ErrorResponse{Error: "Dinosaur not found."})
		return
	}

	if err := dbConn.Delete(&dinosaur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to remove dinosaur."})
		return
	}

	c.JSON(http.StatusOK, apimodels.RemoveDinosaurResponse{})
}

func canBeMovedToCage(c *gin.Context, dbConn *gorm.DB, dinosaur dbmodels.Dinosaur, cageID uint) bool {
	var cage dbmodels.Cage
	if err := dbConn.Preload("Dinosaurs").First(&cage, cageID).Error; err != nil {
		c.JSON(http.StatusNotFound, apimodels.ErrorResponse{Error: "Cage not found."})
		return false
	}

	if cage.Capacity == len(cage.Dinosaurs) {
		c.JSON(http.StatusConflict, apimodels.ErrorResponse{Error: "Dinosaur cannot be placed in cage that is already full."})
		return false
	}

	if cage.PowerStatus != string(apimodels.Active) {
		c.JSON(http.StatusConflict, apimodels.ErrorResponse{Error: "Dinosaur cannot be placed in cage that has no power."})
		return false
	}

	for _, dinosaurInCage := range cage.Dinosaurs {
		if dinosaur.Type == string(apimodels.Herbivore) && dinosaurInCage.Type == string(apimodels.Carnivore) {
			c.JSON(http.StatusConflict, apimodels.ErrorResponse{Error: "Herbivore Dinosaur cannot be placed in cage with Carnivores."})
			return false
		}
		if dinosaur.Type == string(apimodels.Carnivore) && dinosaur.Species != dinosaurInCage.Species {
			c.JSON(http.StatusConflict, apimodels.ErrorResponse{Error: "Carnivore Dinosaur cannot be placed in cage with any other species."})
			return false
		}
	}
	return true
}
