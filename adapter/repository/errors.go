package repository

import "errors"

const (
	errDupEntry = 1062
)

var (
	errDatabase = errors.New("database error")
)
