package application

import "github.com/juankohler/crypto-bot/libs/go/errors"

var (
	ErrInternal = errors.Define("INTERNAL")
	ErrNotFound = errors.Define("NOT_FOUND")
)
