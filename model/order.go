package model

type Order struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Item    []Item `gorm:"many2many:order_items;"`
	ItemIDs []uint `gorm:"-"`
	Total   float64
	Status  string
}
