package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id" redis:"user_id" validate:"omitempty"`
	FirstName string    `json:"first_name" db:"first_name" redis:"first_name" validate:"omitempty,lte=30"`
	LastName  string    `json:"last_name" db:"last_name" redis:"last_name" validate:"omitempty,lte=30"`
	Email     string    `json:"email,omitempty" db:"email" redis:"email" validate:"omitempty,required,lte=60,email"`
	Password  string    `json:"password,omitempty" db:"password" redis:"password" validate:"omitempty,required,gte=6"`
	Nickname  string    `json:"nickname" db:"nickname" redis:"nickname" validate:"omitempty,required,lte=30"`
}

func (u *User) HashPassword() error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPass)
	return nil
}

func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}

	return nil
}

func (u *User) SanitizePassword() {
	u.Password = ""
}

func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(u.Email)
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	return nil
}

func (u *User) PrepareUpdate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	return nil
}

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
