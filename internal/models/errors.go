package models

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: No matching record found")
	ErrInValidCredentials = errors.New("models: invalid Credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)
