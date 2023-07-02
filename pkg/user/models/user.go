package models

import (
	"gorm.io/datatypes"
)

// TableName overrides the table name used by User to `profiles`
func (User) TableName() string {
	return "users"
}

type User struct {
	ID        uint           `gorm:"primaryKey"`
	Email     string         `form:"email" json:"email"`
	Name      string         `form:"name" json:"name"`
	Password  string         `form:"password" json:"password"`
	CreatedAt datatypes.Time `gorm:"->"`
	UpdatedAt datatypes.Time `gorm:"->"`
}
