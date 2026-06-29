package users

import (
	"time"

	"github.com/Rowkash/go-gin-auth/internal/common/types"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"unique;not null;index"`
	Password  string         `json:"-" gorm:"not null;"`
	Role      types.UserRole `json:"role" gorm:"type:varchar(20);default:'USER'"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt time.Time      `json:"-" gorm:"column:updatedAt"`
}

type UserInsertInput struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}
