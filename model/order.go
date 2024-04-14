package model

type Order struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	// Item   []Item
	Total  float64
	Status string
}
