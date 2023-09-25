package handlers

import (
	"net/http"
	"strconv"

	"pp-jurassic-park-api/internal/db"
	transform "pp-jurassic-park-api/internal/models"
	apimodels "pp-jurassic-park-api/internal/models/api"
	dbmodels "pp-jurassic-park-api/internal/models/db"

	"github.com/gin-gonic/gin"
)

func GetCage(c *gin.Context) {
	idParam := c.Param("id")
	cageID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cage ID"})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database."})
		return
	}

	var cage dbmodels.Cage
	if err := dbConn.Preload("Dinosaurs").First(&cage, cageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cage not found."})
		return
	}

	c.JSON(http.StatusOK, apimodels.GetCageResponse{Cage: transform.CageToApi(cage)})
}

func GetCages(c *gin.Context) {
	var req apimodels.GetCagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database."})
		return
	}

	var cages []dbmodels.Cage
	if len(req.FilteredPowerStatuses) > 0 {
		if err := dbConn.Preload("Dinosaurs").Where("power_status IN ?", req.FilteredPowerStatuses).Find(&cages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cages with filtered statuses."})
			return
		}
	} else {
		if err := dbConn.Preload("Dinosaurs").Find(&cages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cages."})
			return
		}
	}

	c.JSON(http.StatusOK, apimodels.GetCagesResponse{Cages: transform.CagesToApi(cages)})
}

func CreateCage(c *gin.Context) {
	var req apimodels.CreateCageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.PowerStatus != apimodels.Active && req.PowerStatus != apimodels.Down {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid power status"})
		return
	}

	if req.Capacity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Capacity should be greater than 0"})
		return
	}

	cage := dbmodels.Cage{
		Capacity:    req.Capacity,
		PowerStatus: string(req.PowerStatus),
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database."})
		return
	}

	if err := dbConn.Create(&cage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cage"})
		return
	}

	c.JSON(http.StatusOK, apimodels.CreateCageResponse{Cage: transform.CageToApi(cage)})
}

func UpdateCagePowerStatus(c *gin.Context) {
	idParam := c.Param("id")
	cageID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cage ID"})
		return
	}

	var req apimodels.UpdateCagePowerStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.PowerStatus != apimodels.Active && req.PowerStatus != apimodels.Down {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid power status"})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database."})
		return
	}

	var cage dbmodels.Cage
	if err := dbConn.First(&cage, cageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cage not found."})
		return
	}

	if cage.PowerStatus != string(req.PowerStatus) {
		cage.PowerStatus = string(req.PowerStatus)
		if err := dbConn.Save(&cage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cage power status."})
			return
		}
	}

	c.JSON(http.StatusOK, apimodels.UpdateCagePowerStatusResponse{Cage: transform.CageToApi(cage)})
}

func DeleteCage(c *gin.Context) {
	idParam := c.Param("id")
	cageID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cage ID"})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database."})
		return
	}

	var cage dbmodels.Cage
	if err := dbConn.First(&cage, cageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cage not found."})
		return
	}

	if len(cage.Dinosaurs) > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete cage with dinosaurs inside."})
		return
	}

	if err := dbConn.Delete(&cage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cage."})
		return
	}

	c.JSON(http.StatusOK, apimodels.DeleteCageResponse{})
}
