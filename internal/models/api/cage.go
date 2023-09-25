package apimodels

type PowerStatus string

const (
	Active PowerStatus = "ACTIVE"
	Down   PowerStatus = "DOWN"
)

type Cage struct {
	ID           uint        `json:"id"`
	Capacity     uint        `json:"capacity"`
	CurrentCount int         `json:"current_count"`
	PowerStatus  PowerStatus `json:"power_status"`
	Dinosaurs    []Dinosaur  `json:"dinosaurs"`
}

type CreateCageRequest struct {
	Capacity    uint        `json:"capacity"`
	PowerStatus PowerStatus `json:"power_status"`
}
type CreateCageResponse struct {
	Cage Cage `json:"cage"`
}

type UpdateCagePowerStatusRequest struct {
	PowerStatus PowerStatus `json:"power_status"`
}
type UpdateCagePowerStatusResponse struct {
	Cage Cage `json:"cage"`
}

type DeleteCageRequest struct {
}
type DeleteCageResponse struct {
}

type GetCagesRequest struct {
	FilteredPowerStatuses []PowerStatus `json:"filtered_power_status,omitempty"`
}
type GetCagesResponse struct {
	Cages []Cage `json:"cages"`
}

type GetCageRequest struct {
}
type GetCageResponse struct {
	Cage Cage `json:"cage"`
}
