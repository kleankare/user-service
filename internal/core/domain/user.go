package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string  `gorm:"size:256;not null;unique" json:"username"`
	Password string  `gorm:"size:256;not null" json:"password"`
	Role     []*Role `gorm:"many2many:user_roles;"`
}
