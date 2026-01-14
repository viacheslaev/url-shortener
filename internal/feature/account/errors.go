package account

import "errors"

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrEmailInvalidFormat    = errors.New("incorrect email format")
	ErrPasswordInvalidFormat = errors.New("password must be at least 6 characters")
	ErrAccountNotFound       = errors.New("account not found")
)
