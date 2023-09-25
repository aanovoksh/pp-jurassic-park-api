package dbmodels

type Cage struct {
	ID          uint       `gorm:"primaryKey;autoIncrement"`
	Capacity    uint       `gorm:"not null"`
	PowerStatus string     `gorm:"not null"`
	Dinosaurs   []Dinosaur `gorm:"foreignKey:CageID"`
}
