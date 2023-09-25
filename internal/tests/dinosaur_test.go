package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	apimodels "pp-jurassic-park-api/internal/api/models"
	"pp-jurassic-park-api/internal/db"
	dbmodels "pp-jurassic-park-api/internal/db/models"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDinosaur(t *testing.T) {
	t.Run("Dinosaur found", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/dinosaurs/"+strconv.FormatUint(uint64(tyrannosaurus.ID), 10), nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var getResponse apimodels.GetDinosaurResponse
		json.Unmarshal(response.Body.Bytes(), &getResponse)

		assertDinosaur(t, getResponse.Dinosaur, tyrannosaurus.Name, apimodels.Species(tyrannosaurus.Species), apimodels.DinosaurType(tyrannosaurus.Type), tyrannosaurus.CageID)
	})

	t.Run("Invalid dinosaur ID", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/dinosaurs/invalidID", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Dinosaur not found", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/dinosaurs/123456", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestGetDinosaurs(t *testing.T) {
	router := setupRouter()

	t.Run("Retrieve all dinosaurs", func(t *testing.T) {
		data, _ := json.Marshal(apimodels.GetDinosaursRequest{FilteredSpecies: []apimodels.Species{}})
		request, _ := http.NewRequest(http.MethodGet, "/dinosaurs", bytes.NewReader(data))

		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var getResponse apimodels.GetDinosaursResponse
		json.Unmarshal(response.Body.Bytes(), &getResponse)

		assert.Len(t, getResponse.Dinosaurs, 5)
	})

	t.Run("Retrieve dinosaurs with filtered species", func(t *testing.T) {
		data, _ := json.Marshal(apimodels.GetDinosaursRequest{FilteredSpecies: []apimodels.Species{apimodels.Tyrannosaurus}})
		request, _ := http.NewRequest(http.MethodGet, "/dinosaurs", bytes.NewReader(data))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var getResponse apimodels.GetDinosaursResponse
		json.Unmarshal(response.Body.Bytes(), &getResponse)

		assert.Len(t, getResponse.Dinosaurs, 1)
	})

	t.Run("Invalid input", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/dinosaurs", strings.NewReader("invalid"))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestCreateDinosaur(t *testing.T) {
	t.Run("Successful dinosaur creation", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"name": "Barry",
			"species": "Brachiosaurus",
			"cage_id": %d
		}`, activeCage.ID)
		request, _ := http.NewRequest(http.MethodPost, "/dinosaurs", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var createResponse apimodels.AddDinosaurResponse
		json.Unmarshal(response.Body.Bytes(), &createResponse)
		dinosaurIDsToCleanup = append(dinosaurIDsToCleanup, createResponse.Dinosaur.ID)
		assertDinosaur(t, createResponse.Dinosaur, "Barry", apimodels.Brachiosaurus, apimodels.Herbivore, activeCage.ID)
	})

	t.Run("Invalid input", func(t *testing.T) {
		payload := `{
			"InvalidJson"
		}`
		request, _ := http.NewRequest(http.MethodPost, "/dinosaurs", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid name", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"name": "  ",
			"species": "Brachiosaurus",
			"cage_id": %d
		}`, activeCage.ID)
		request, _ := http.NewRequest(http.MethodPost, "/dinosaurs", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid species", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"name": "Minni",
			"species": "Invalid",
			"cage_id": %d
		}`, activeCage.ID)
		request, _ := http.NewRequest(http.MethodPost, "/dinosaurs", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid cage: power status is down", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"name": "Barry",
			"species": "Brachiosaurus",
			"cage_id": %d
		}`, downCage.ID)
		request, _ := http.NewRequest(http.MethodPost, "/dinosaurs", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("Invalid cage: capacity", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"name": "Barry",
			"species": "Brachiosaurus",
			"cage_id": %d
		}`, activeCage.ID)
		request, _ := http.NewRequest(http.MethodPost, "/dinosaurs", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("Invalid cage: herbivores cannot share cage with carnivores", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"name": "Barry",
			"species": "Brachiosaurus",
			"cage_id": %d
		}`, cageWithTyrannosaurus.ID)
		request, _ := http.NewRequest(http.MethodPost, "/dinosaurs", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("Invalid cage: carnivore species cannot share cage with anyone", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"name": "Terry",
			"species": "Tyrannosaurus",
			"cage_id": %d
		}`, cageWithSpinosaurus.ID)
		request, _ := http.NewRequest(http.MethodPost, "/dinosaurs", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
	})
}

func TestMoveDinosaur(t *testing.T) {
	t.Run("Successful move of a dinosaur", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"cage_id": %d
		}`, cageWithStegosaurus.ID)
		request, _ := http.NewRequest(http.MethodPatch, "/dinosaurs/"+strconv.FormatUint(uint64(triceratops.ID), 10), bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var createResponse apimodels.MoveDinosaurResponse
		json.Unmarshal(response.Body.Bytes(), &createResponse)
		assertDinosaur(t, createResponse.Dinosaur, triceratops.Name, apimodels.Species(triceratops.Species), apimodels.DinosaurType(triceratops.Type), cageWithStegosaurus.ID)
	})

	t.Run("Invalid dinosaur ID", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"cage_id": %d
		}`, cageWithStegosaurus.ID)
		request, _ := http.NewRequest(http.MethodPatch, "/dinosaurs/invalidID", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid input", func(t *testing.T) {
		payload := `{
			"InvalidJson"
		}`
		request, _ := http.NewRequest(http.MethodPatch, "/dinosaurs/"+strconv.FormatUint(uint64(triceratops.ID), 10), bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid cage id", func(t *testing.T) {
		payload := `{
			"cage_id": "InvalidCageId"
		}`
		request, _ := http.NewRequest(http.MethodPatch, "/dinosaurs/"+strconv.FormatUint(uint64(triceratops.ID), 10), bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Dinosaurs not found", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"cage_id": %d
		}`, cageWithStegosaurus.ID)
		request, _ := http.NewRequest(http.MethodPatch, "/dinosaurs/123456", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("Invalid cage: power status is down", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"cage_id": %d
		}`, downCage.ID)
		request, _ := http.NewRequest(http.MethodPatch, "/dinosaurs/"+strconv.FormatUint(uint64(triceratops.ID), 10), bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("Invalid cage: capacity", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"cage_id": %d
		}`, activeCage.ID)
		request, _ := http.NewRequest(http.MethodPatch, "/dinosaurs/"+strconv.FormatUint(uint64(triceratops.ID), 10), bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("Invalid cage: herbivores cannot share cage with carnivores", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"cage_id": %d
		}`, cageWithTyrannosaurus.ID)
		request, _ := http.NewRequest(http.MethodPatch, "/dinosaurs/"+strconv.FormatUint(uint64(triceratops.ID), 10), bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("Invalid cage: carnivore species cannot share cage with anyone", func(t *testing.T) {
		payload := fmt.Sprintf(`{
			"cage_id": %d
		}`, cageWithStegosaurus.ID)
		request, _ := http.NewRequest(http.MethodPatch, "/dinosaurs/"+strconv.FormatUint(uint64(tyrannosaurus.ID), 10), bytes.NewBufferString(payload))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
	})
}

func TestRemoveDinosaur(t *testing.T) {
	t.Run("Remove existing dinosaur", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/dinosaurs/"+strconv.FormatUint(uint64(cageToBeUpdated.ID), 10), nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Invalid dinosaur ID", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/dinosaurs/invalidID", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Removing non-existent dinosaur", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/dinosaurs/123456", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func CreateTestDinosaur(name string, species apimodels.Species, dinosaurType apimodels.DinosaurType, cageID uint) dbmodels.Dinosaur {
	dbConn, _ := db.Connect()

	dinosaur := dbmodels.Dinosaur{
		Name:    name,
		Species: string(species),
		Type:    string(dinosaurType),
		CageID:  cageID,
	}
	dbConn.Create(&dinosaur)
	dinosaurIDsToCleanup = append(dinosaurIDsToCleanup, dinosaur.ID)
	return dinosaur
}

func DeleteTestDinosaurs(ids []uint) {
	dbConn, _ := db.Connect()
	dbConn.Where("id IN (?)", ids).Delete(&dbmodels.Dinosaur{})
}

func assertDinosaur(t *testing.T, dinosaur apimodels.Dinosaur, name string, species apimodels.Species, dinosaurType apimodels.DinosaurType, cageID uint) {
	assert.Equal(t, name, dinosaur.Name)
	assert.Equal(t, species, dinosaur.Species)
	assert.Equal(t, dinosaurType, dinosaur.Type)
	assert.Equal(t, cageID, dinosaur.CageID)
}
