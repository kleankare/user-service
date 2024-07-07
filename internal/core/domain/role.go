package domain

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name string  `gorm:"size:256;not null;unique" json:"name"`
	User []*User `gorm:"many2many:user_roles;"`
}
