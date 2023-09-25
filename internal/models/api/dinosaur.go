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

type Dinosaur struct {
	ID      uint         `json:"id"`
	Name    string       `json:"name"`
	Species Species      `json:"species"`
	Type    DinosaurType `json:"type"`
	CageID  uint         `json:"cage_id"`
}
