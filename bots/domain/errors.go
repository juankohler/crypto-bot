package domain

import (
	"github.com/juankohler/crypto-bot/libs/go/errors"
)

var (
	ErrInvalid  = errors.Define("INVALID")
	ErrNotFound = errors.Define("NOT_FOUND")
	ErrInternal = errors.Define("INTERNAL")
)
