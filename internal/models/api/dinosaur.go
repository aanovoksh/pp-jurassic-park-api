package apimodels

type DinosaurType string

const (
	Herbivore DinosaurType = "HERBIVORE"
	Carnivore DinosaurType = "CARNIVORE"
)

type Species string

const (
	Tyrannosaurus Species = "Tyrannosaurus"
	Velociraptor  Species = "Velociraptor"
	Spinosaurus   Species = "Spinosaurus"
	Megalosaurus  Species = "Megalosaurus"
	Brachiosaurus Species = "Brachiosaurus"
	Stegosaurus   Species = "Stegosaurus"
	Ankylosaurus  Species = "Ankylosaurus"
	Triceratops   Species = "Triceratops"
)

var speciesToTypeLookup = map[Species]DinosaurType{
	Tyrannosaurus: Carnivore,
	Velociraptor:  Carnivore,
	Spinosaurus:   Carnivore,
	Megalosaurus:  Carnivore,
	Brachiosaurus: Herbivore,
	Stegosaurus:   Herbivore,
	Ankylosaurus:  Herbivore,
	Triceratops:   Herbivore,
}

func LookupSpeciesType(species string) (bool, Species, DinosaurType) {
	dinosaurType, exists := speciesToTypeLookup[Species(species)]
	return exists, Species(species), dinosaurType
}

type Dinosaur struct {
	ID      uint         `json:"id"`
	Name    string       `json:"name"`
	Species Species      `json:"species"`
	Type    DinosaurType `json:"type"`
	CageID  uint         `json:"cage_id"`
}

type AddDinosaurRequest struct {
	Name    string `json:"name"`
	Species string `json:"species"`
	CageID  uint   `json:"cage_id"`
}
type AddDinosaurResponse struct {
	Dinosaur Dinosaur `json:"dinosaur"`
}

type MoveDinosaurRequest struct {
	CageID uint `json:"cage_id"`
}
type MoveDinosaurResponse struct {
	Dinosaur Dinosaur `json:"dinosaur"`
}

type RemoveDinosaurRequest struct {
}
type RemoveDinosaurResponse struct {
}

type GetDinosaurRequest struct {
}
type GetDinosaurResponse struct {
	Dinosaur Dinosaur `json:"dinosaur"`
}

type GetDinosaursRequest struct {
	FilteredSpecies []Species `json:"filtered_species,omitempty"`
}
type GetDinosaursResponse struct {
	Dinosaurs []Dinosaur `json:"dinosaurs"`
}
