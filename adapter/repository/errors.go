package repository

import "errors"

const (
	ERR_DUP_ENTRY = 1062
)

var (
	ErrDatabase = errors.New("database error")
)
