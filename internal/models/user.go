// Package models describes database tables
package models

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	Regular UserRole = "regular"
	Admin   UserRole = "admin"
)

var Roles = []UserRole{
	Regular,
	Admin,
}

type User struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Name     string    `json:"name" gorm:"not null"`
	Email    string    `json:"email" gorm:"not null;uniqueIndex"`
	Password string    `json:"-"`
	Role     UserRole  `json:"role" gorm:"default:'regular'"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (role *UserRole) Scan(value interface{}) error {
	strValue, ok := value.(string)

	if !ok {
		return errors.New("invalid value for UserRole")
	}

	*role = UserRole(strValue)

	if !role.IsValid() {
		return fmt.Errorf("invalid role value: %s", strValue)
	}

	return nil
}

func (role UserRole) IsValid() bool {
	return slices.Contains(Roles, role)
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()

	if !user.Role.IsValid() {
		return fmt.Errorf("invalid role: %s", user.Role)
	}

	if user.Role == Admin {
		return errors.New("cannot set admin role")
	}

	user.HashPassword(user.Password)

	return nil
}

func (user *User) HashPassword(password string) error {
	hashedPasswordbytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("unable to hash password")
	}

	user.Password = string(hashedPasswordbytes)

	return nil
}

func (user *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (user *User) HasRole(role UserRole) bool {
	return user.Role == role
}

func (user *User) IsAdmin() bool {
	return user.Role == Admin
}
