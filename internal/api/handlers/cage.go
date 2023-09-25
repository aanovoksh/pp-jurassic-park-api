package handlers

import (
	"net/http"
	"strconv"

	apimodels "pp-jurassic-park-api/internal/api/models"
	transform "pp-jurassic-park-api/internal/api/transform"
	db "pp-jurassic-park-api/internal/db"
	dbmodels "pp-jurassic-park-api/internal/db/models"

	"github.com/gin-gonic/gin"
)

// GetCage returns single cage for the requested id.
// Used to retrieve data around single cage and its habitants at the Jurassic Park.
func GetCage(c *gin.Context) {
	idParam := c.Param("id")
	cageID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid cage ID."})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to connect to database."})
		return
	}

	var cage dbmodels.Cage
	if err := dbConn.Preload("Dinosaurs").First(&cage, cageID).Error; err != nil {
		c.JSON(http.StatusNotFound, apimodels.ErrorResponse{Error: "Cage not found."})
		return
	}

	c.JSON(http.StatusOK, apimodels.GetCageResponse{Cage: transform.CageToApi(cage)})
}

// GetDinosaurs returns all dinosaurs matching provided filters.
// Used to retrieve data around all cages and their habitants at the Jurassic Park.
func GetCages(c *gin.Context) {
	var req apimodels.GetCagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid input."})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to connect to database."})
		return
	}

	var cages []dbmodels.Cage
	if len(req.FilteredPowerStatuses) > 0 {
		if err := dbConn.Preload("Dinosaurs").Where("power_status IN ?", req.FilteredPowerStatuses).Find(&cages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to retrieve cages with filtered statuses."})
			return
		}
	} else {
		if err := dbConn.Preload("Dinosaurs").Find(&cages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to retrieve cages."})
			return
		}
	}

	c.JSON(http.StatusOK, apimodels.GetCagesResponse{Cages: transform.CagesToApi(cages)})
}

// UpdateCagePowerStatus creates a new cage.
// Used to register a new cage at the Jurassic Park.
func CreateCage(c *gin.Context) {
	var req apimodels.CreateCageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid input."})
		return
	}

	if req.PowerStatus != apimodels.Active && req.PowerStatus != apimodels.Down {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid power status."})
		return
	}

	if req.Capacity <= 0 {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Capacity should be greater than 0."})
		return
	}

	cage := dbmodels.Cage{
		Capacity:    req.Capacity,
		PowerStatus: string(req.PowerStatus),
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to connect to database."})
		return
	}

	if err := dbConn.Create(&cage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to create cage."})
		return
	}

	c.JSON(http.StatusOK, apimodels.CreateCageResponse{Cage: transform.CageToApi(cage)})
}

// UpdateCagePowerStatus sets power status for a given cage.
// Used to control power of cages at the Jurassic Park.
func UpdateCagePowerStatus(c *gin.Context) {
	idParam := c.Param("id")
	cageID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid cage ID."})
		return
	}

	var req apimodels.UpdateCagePowerStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid input."})
		return
	}

	if req.PowerStatus != apimodels.Active && req.PowerStatus != apimodels.Down {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid power status."})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to connect to database."})
		return
	}

	var cage dbmodels.Cage
	if err := dbConn.First(&cage, cageID).Error; err != nil {
		c.JSON(http.StatusNotFound, apimodels.ErrorResponse{Error: "Cage not found."})
		return
	}

	if cage.PowerStatus != string(req.PowerStatus) {
		cage.PowerStatus = string(req.PowerStatus)
		if err := dbConn.Save(&cage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to update cage power status."})
			return
		}
	}

	c.JSON(http.StatusOK, apimodels.UpdateCagePowerStatusResponse{Cage: transform.CageToApi(cage)})
}

// DeleteCage deletes the cage.
// Used to remove all no longer needed cages at the Jurassic Park.
func DeleteCage(c *gin.Context) {
	idParam := c.Param("id")
	cageID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodels.ErrorResponse{Error: "Invalid cage ID."})
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to connect to database."})
		return
	}

	var cage dbmodels.Cage
	if err := dbConn.Preload("Dinosaurs").First(&cage, cageID).Error; err != nil {
		c.JSON(http.StatusNotFound, apimodels.ErrorResponse{Error: "Cage not found."})
		return
	}

	if len(cage.Dinosaurs) > 0 {
		c.JSON(http.StatusConflict, apimodels.ErrorResponse{Error: "Cannot delete cage with dinosaurs inside."})
		return
	}

	if err := dbConn.Delete(&cage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, apimodels.ErrorResponse{Error: "Failed to delete cage."})
		return
	}

	c.JSON(http.StatusOK, apimodels.DeleteCageResponse{})
}
