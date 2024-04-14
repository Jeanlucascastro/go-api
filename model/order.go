package model

type Order struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Item   []Item `gorm:"many2many:order_items;"`
	Total  float64
	Status string
}
