package models

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jaevor/go-nanoid"
	"github.com/juankohler/crypto-bot/libs/go/errors"
)

var (
	ErrInvalidID = errors.Define("id.invalid")
)

type ID string

func NewID(id string) (ID, error) {
	id = strings.TrimSpace(id)

	if len(id) < 5 {
		return "", errors.New(ErrInvalidID, "short id", errors.WithMetadata("id", id))
	}

	return ID(id), nil
}

var alphabetNanoID = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateNanoID(lenght int) (ID, error) {
	gen, err := nanoid.CustomASCII(alphabetNanoID, lenght)
	if err != nil {
		return "", err
	}
	return ID(gen()), nil
}

func GenerateUUID() ID {
	return ID(uuid.New().String())
}

func GenerateSlug(str string) (ID, error) {
	s := slug.Make(str)
	return NewID(s)
}

func (id ID) String() string {
	return string(id)
}

// Serialization
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	bID, err := NewID(s)
	if err != nil {
		return err
	}

	*id = bID

	return nil
}
