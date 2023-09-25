package dbmodels

type Dinosaur struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Name    string `gorm:"not null"`
	Species string `gorm:"not null"`
	Type    string `gorm:"not null"`
	CageID  uint   `gorm:"not null"`
}
