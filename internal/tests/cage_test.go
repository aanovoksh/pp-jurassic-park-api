package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	apimodels "pp-jurassic-park-api/internal/api/models"
	"pp-jurassic-park-api/internal/db"
	dbmodels "pp-jurassic-park-api/internal/db/models"

	"github.com/stretchr/testify/assert"
)

func TestGetCage(t *testing.T) {
	t.Run("Cage found", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/cages/"+strconv.FormatUint(uint64(activeCage.ID), 10), nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Invalid cage ID", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/cages/invalidID", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Cage not found", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/cages/123456", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestGetCages(t *testing.T) {
	t.Run("Retrieve all cages", func(t *testing.T) {
		data, _ := json.Marshal(apimodels.GetCagesRequest{FilteredPowerStatuses: []apimodels.PowerStatus{}})
		request, _ := http.NewRequest(http.MethodGet, "/cages", bytes.NewReader(data))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var getResponse apimodels.GetCagesResponse
		json.Unmarshal(response.Body.Bytes(), &getResponse)

		assert.Len(t, getResponse.Cages, 8)
	})

	t.Run("Retrieve cages with specific power status", func(t *testing.T) {
		data, _ := json.Marshal(apimodels.GetCagesRequest{FilteredPowerStatuses: []apimodels.PowerStatus{apimodels.Active}})
		request, _ := http.NewRequest(http.MethodGet, "/cages", bytes.NewReader(data))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var getResponse apimodels.GetCagesResponse
		json.Unmarshal(response.Body.Bytes(), &getResponse)

		assert.Len(t, getResponse.Cages, 6)
		assert.Equal(t, apimodels.Active, getResponse.Cages[0].PowerStatus)
	})

	t.Run("Invalid input", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/cages", strings.NewReader("invalid"))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestCreateCage(t *testing.T) {
	t.Run("Successful cage creation", func(t *testing.T) {
		payload := `{
			"power_status": "ACTIVE",
			"capacity": 5
		}`
		request, _ := http.NewRequest(http.MethodPost, "/cages", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var createResponse apimodels.CreateCageResponse
		json.Unmarshal(response.Body.Bytes(), &createResponse)
		cageIDsToCleanup = append(cageIDsToCleanup, createResponse.Cage.ID)
		assertCage(t, createResponse.Cage, 5, apimodels.Active, 0)
	})

	t.Run("Invalid input", func(t *testing.T) {
		payload := `{
			"InvalidJson"
		}`
		request, _ := http.NewRequest(http.MethodPost, "/cages", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid power status", func(t *testing.T) {
		payload := `{
			"power_status": "InvalidStatus",
			"capacity": 5
		}`
		request, _ := http.NewRequest(http.MethodPost, "/cages", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid capacity", func(t *testing.T) {
		payload := `{
			"power_status": "ACTIVE",
			"capacity": 0
		}`
		request, _ := http.NewRequest(http.MethodPost, "/cages", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestUpdateCagePowerStatus(t *testing.T) {
	t.Run("Successful power status update", func(t *testing.T) {
		payload := `{
			"power_status": "DOWN"
		}`
		request, _ := http.NewRequest(http.MethodPatch, "/cages/"+strconv.FormatUint(uint64(cageToBeUpdated.ID), 10), bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var updateResponse apimodels.UpdateCagePowerStatusResponse
		json.Unmarshal(response.Body.Bytes(), &updateResponse)
		assertCage(t, updateResponse.Cage, 3, apimodels.Down, 0)
	})

	t.Run("Invalid cage ID", func(t *testing.T) {
		payload := `{
			"power_status": "ACTIVE"
		}`
		request, _ := http.NewRequest(http.MethodPatch, "/cages/invalidID", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid input", func(t *testing.T) {
		payload := `{
			"InvalidJson"
		}`
		request, _ := http.NewRequest(http.MethodPatch, "/cages/"+strconv.FormatUint(uint64(cageToBeUpdated.ID), 10), bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid power status", func(t *testing.T) {
		payload := `{
			"power_status": "InvalidStatus"
		}`
		request, _ := http.NewRequest(http.MethodPatch, "/cages/1", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Cage not found", func(t *testing.T) {
		payload := `{
			"power_status": "ACTIVE"
		}`
		request, _ := http.NewRequest(http.MethodPatch, "/cages/123456", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestDeleteCage(t *testing.T) {
	t.Run("Successful cage deletion", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/cages/"+strconv.FormatUint(uint64(cageToBeRemoved.ID), 10), nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Invalid cage ID", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/cages/invalidID", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Cage not found", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/cages/123456", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("Cage containing dinosaurs", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/cages/"+strconv.FormatUint(uint64(cageWithTyrannosaurus.ID), 10), nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
	})
}

func CreateTestCage(capacity int, powerStatus apimodels.PowerStatus) dbmodels.Cage {
	dbConn, _ := db.Connect()

	cage := dbmodels.Cage{
		Capacity:    capacity,
		PowerStatus: string(powerStatus),
	}
	dbConn.Create(&cage)
	cageIDsToCleanup = append(cageIDsToCleanup, cage.ID)
	return cage
}

func DeleteTestCages(ids []uint) {
	dbConn, _ := db.Connect()
	dbConn.Where("id IN (?)", ids).Delete(&dbmodels.Cage{})
}

func assertCage(t *testing.T, cage apimodels.Cage, capacity int, powerStatus apimodels.PowerStatus, currentCount int) {
	assert.Equal(t, capacity, cage.Capacity)
	assert.Equal(t, powerStatus, cage.PowerStatus)
	assert.Equal(t, currentCount, cage.CurrentCount)
	assert.Greater(t, cage.ID, uint(0))
}
