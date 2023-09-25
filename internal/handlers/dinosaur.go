package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"pp-jurassic-park-api/internal/db"
	transform "pp-jurassic-park-api/internal/models"
	apimodels "pp-jurassic-park-api/internal/models/api"
	dbmodels "pp-jurassic-park-api/internal/models/db"

	"github.com/gin-gonic/gin"
)

func GetDinosaur(c *gin.Context) {
	idParam := c.Param("id")
	dinosaurID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dinosaur ID"})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database."})
		return
	}

	var dinosaur dbmodels.Dinosaur
	if err := dbConn.First(&dinosaur, dinosaurID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dinosaur not found."})
		return
	}

	c.JSON(http.StatusOK, apimodels.GetDinosaurResponse{Dinosaur: transform.DinosaurToApi(dinosaur)})
}

func GetDinosaurs(c *gin.Context) {
	var req apimodels.GetDinosaursRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database."})
		return
	}

	var dinosaurs []dbmodels.Dinosaur
	if len(req.FilteredSpecies) > 0 {
		if err := dbConn.Where("species IN ?", req.FilteredSpecies).Find(&dinosaurs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dinosaurs with filtered species."})
			return
		}
	} else {
		if err := dbConn.Find(&dinosaurs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dinosaurs."})
			return
		}
	}

	c.JSON(http.StatusOK, apimodels.GetDinosaursResponse{Dinosaurs: transform.DinosaursToApi(dinosaurs)})
}

func AddDinosaur(c *gin.Context) {
	var req apimodels.AddDinosaurRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database."})
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name. Name cannot be blank."})
		return
	}

	knownSpecies, species, dinosaurType := apimodels.LookupSpeciesType(req.Species)
	if !knownSpecies {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown species."})
		return
	}

	var cage dbmodels.Cage
	if err := dbConn.Preload("Dinosaurs").First(&cage, req.CageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cage not found."})
		return
	}

	if cage.Capacity == uint(len(cage.Dinosaurs)) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dinosaur cannot be placed in cage that is already full."})
		return
	}

	if cage.PowerStatus != string(apimodels.Active) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dinosaur cannot be placed in cage that has no power."})
		return
	}

	for _, dinosaurInCage := range cage.Dinosaurs {
		if dinosaurType == apimodels.Herbivore && dinosaurInCage.Type == string(apimodels.Carnivore) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Herbivore Dinosaur cannot be placed in cage with Carnivores."})
			return
		}
		if dinosaurType == apimodels.Carnivore && dinosaurInCage.Species != string(species) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carnivore Dinosaur cannot be placed in cage with any other species."})
			return
		}
	}

	dinosaur := dbmodels.Dinosaur{
		Name:    req.Name,
		Species: string(species),
		Type:    string(dinosaurType),
		CageID:  req.CageID,
	}
	if err := dbConn.Create(&dinosaur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add dinosaur"})
		return
	}

	c.JSON(http.StatusOK, apimodels.AddDinosaurResponse{Dinosaur: transform.DinosaurToApi(dinosaur)})
}

func MoveDinosaur(c *gin.Context) {
	idParam := c.Param("id")
	dinosaurID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dinosaur ID"})
		return
	}

	var req apimodels.MoveDinosaurRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database."})
		return
	}

	var dinosaur dbmodels.Dinosaur
	if err := dbConn.First(&dinosaur, dinosaurID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dinosaur not found."})
		return
	}

	var cage dbmodels.Cage
	if err := dbConn.Preload("Dinosaurs").First(&cage, req.CageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cage not found."})
		return
	}

	if cage.Capacity == uint(len(cage.Dinosaurs)) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dinosaur cannot be placed in cage that is already full."})
		return
	}

	if cage.PowerStatus != string(apimodels.Active) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dinosaur cannot be placed in cage that has no power."})
		return
	}

	for _, dinosaurInCage := range cage.Dinosaurs {
		if dinosaur.Type == string(apimodels.Herbivore) && dinosaurInCage.Type == string(apimodels.Carnivore) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Herbivore Dinosaur cannot be placed in cage with Carnivores."})
			return
		}
		if dinosaur.Type == string(apimodels.Carnivore) && dinosaurInCage.Species != dinosaur.Species {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carnivore Dinosaur cannot be placed in cage with any other species."})
			return
		}
	}

	if dinosaur.CageID != req.CageID {
		dinosaur.CageID = req.CageID
		if err := dbConn.Save(&dinosaur).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move dinosaur."})
			return
		}
	}

	c.JSON(http.StatusOK, apimodels.MoveDinosaurResponse{Dinosaur: transform.DinosaurToApi(dinosaur)})
}

func RemoveDinosaur(c *gin.Context) {
	idParam := c.Param("id")
	dinosaurID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dinosaur ID"})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database."})
		return
	}

	var dinosaur dbmodels.Dinosaur
	if err := dbConn.First(&dinosaur, dinosaurID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dinosaur not found."})
		return
	}

	if err := dbConn.Delete(&dinosaur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove dinosaur."})
		return
	}

	c.JSON(http.StatusOK, apimodels.RemoveDinosaurResponse{})
}
