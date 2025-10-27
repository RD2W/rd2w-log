package model

import (
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Callsign  string    `json:"callsign" db:"callsign"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password_hash"`
	Name      string    `json:"name" db:"name"`
	Location  string    `json:"location" db:"location"`
	Timezone  string    `json:"timezone" db:"timezone"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewUser(callsign, email, password, name, location, timezone string) (*User, error) {
	user := &User{
		ID:        uuid.New(),
		Callsign:  callsign,
		Email:     email,
		Name:      name,
		Location:  location,
		Timezone:  timezone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) SetPassword(password string) error {
	if utf8.RuneCountInString(password) < 8 {
		return ErrPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) Validate() error {
	if u.Callsign == "" {
		return ErrInvalidCallsign
	}

	if u.Email == "" {
		return ErrInvalidEmail
	}

	if u.Password == "" {
		return ErrInvalidPassword
	}

	return nil
}

func (u *User) UpdateProfile(name, location, timezone string) {
	u.Name = name
	u.Location = location
	u.Timezone = timezone
	u.UpdatedAt = time.Now()
}
