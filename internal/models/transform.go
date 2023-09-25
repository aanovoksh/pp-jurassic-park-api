package transform

import (
	apimodels "pp-jurassic-park-api/internal/models/api"
	dbmodels "pp-jurassic-park-api/internal/models/db"
)

func CagesToApi(dbCages []dbmodels.Cage) []apimodels.Cage {
	apiCages := []apimodels.Cage{}
	for _, dbCage := range dbCages {
		apiCages = append(apiCages, CageToApi(dbCage))
	}
	return apiCages
}

func CageToApi(dbCage dbmodels.Cage) apimodels.Cage {
	apiDinosaurs := []apimodels.Dinosaur{}
	for _, dbDino := range dbCage.Dinosaurs {
		apiDinosaurs = append(apiDinosaurs, DinosaurToApi(dbDino))
	}
	return apimodels.Cage{
		ID:           dbCage.ID,
		Capacity:     dbCage.Capacity,
		PowerStatus:  apimodels.PowerStatus(dbCage.PowerStatus),
		CurrentCount: len(apiDinosaurs),
		Dinosaurs:    apiDinosaurs,
	}
}

func DinosaursToApi(dbDinosaurs []dbmodels.Dinosaur) []apimodels.Dinosaur {
	apiDinosaurs := []apimodels.Dinosaur{}
	for _, dbDinosaur := range dbDinosaurs {
		apiDinosaurs = append(apiDinosaurs, DinosaurToApi(dbDinosaur))
	}
	return apiDinosaurs
}

func DinosaurToApi(dbDinosaur dbmodels.Dinosaur) apimodels.Dinosaur {
	return apimodels.Dinosaur{
		ID:      dbDinosaur.ID,
		Name:    dbDinosaur.Name,
		Species: apimodels.Species(dbDinosaur.Species),
		Type:    apimodels.DinosaurType(dbDinosaur.Type),
		CageID:  dbDinosaur.CageID,
	}
}
