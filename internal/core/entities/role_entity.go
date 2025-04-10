package entities

import (
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"gorm.io/gorm"
)

type Role struct {
	ID          utils.ULID `gorm:"column:id;primaryKey;type:varchar(26)"`
	Name        string     `gorm:"unique;not null"`
	Permissions string     `json:"permissions" gorm:"type:json;not null"`

	Users []User `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE"`
}

func (u *Role) TableName() string {
	return "roles"
}

func (u *Role) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = utils.NewULID()
	return nil
}
