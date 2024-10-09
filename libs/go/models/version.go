package models

import (
	"encoding/json"

	"github.com/juankohler/crypto-bot/libs/go/errors"
)

var (
	ErrInvalidVersion = errors.Define("version.invalid")
)

type Version struct {
	updated bool
	Value   int
}

func CreateVersion() Version {
	return Version{
		updated: true,
		Value:   1,
	}
}

func NewVersion(v int) (Version, error) {
	if v < 1 {
		return Version{}, errors.New(
			ErrInvalidVersion,
			"invalid version",
			errors.WithMetadata("version", v),
		)
	}

	return Version{
		Value: v,
	}, nil
}

func (v Version) Update() Version {
	if v.updated {
		return v
	}

	v.updated = true
	v.Value++

	return v
}

// Serialization
func (v Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value)
}

func (v *Version) UnmarshalJSON(data []byte) error {
	var version int
	if err := json.Unmarshal(data, &version); err != nil {
		return err
	}

	bVersion, err := NewVersion(version)
	if err != nil {
		return err
	}

	*v = bVersion

	return nil
}
