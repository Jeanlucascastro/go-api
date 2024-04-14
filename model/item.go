package model

type Item struct {
	ID    uint `gorm:"primaryKey;autoIncrement"`
	Name  string
	Price float64
}
