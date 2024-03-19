package models

type Items struct {
	ID          uint   `gorm:"primaryKey"`
	ItemCode    string `gorm:"not null;type:varchar(191)"`
	Description string `gorm:"not null"`
	Quantity    int
	OrderID     uint
}
