package models

import "time"

type Orders struct {
	ID           uint    `gorm:"primaryKey"`
	CustomerName string  `gorm:"not null;type:varchar(191)"`
	Items        []Items `gorm:"foreignKey:OrderID"` // Menetapkan foreign key secara eksplisit
	OrderedAt    time.Time
}
