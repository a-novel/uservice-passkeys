package dao

import "errors"

var (
	ErrPasskeyNotFound = errors.New("passkey not found")
	ErrInvalidPasskey  = errors.New("invalid passkey")
)
