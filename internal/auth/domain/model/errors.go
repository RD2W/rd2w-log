package model

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCallsign    = errors.New("invalid callsign")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
	ErrSessionNotFound    = errors.New("session not found")
)
