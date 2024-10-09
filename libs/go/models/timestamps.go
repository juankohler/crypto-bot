package models

import (
	"encoding/json"
	"time"

	"github.com/juankohler/crypto-bot/libs/go/errors"
)

var (
	ErrInvalidTimestamps = errors.Define("timestamps.invalid")
)

type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewTimestamps(createdAt time.Time, updatedAt time.Time, deletedAt *time.Time) (Timestamps, error) {
	if createdAt.IsZero() {
		return Timestamps{}, errors.New(
			ErrInvalidTimestamps,
			"invalid created_at",
		)
	}

	if updatedAt.IsZero() {
		return Timestamps{}, errors.New(
			ErrInvalidTimestamps,
			"invalid updated_at",
		)
	}

	if deletedAt != nil && deletedAt.IsZero() {
		return Timestamps{}, errors.New(
			ErrInvalidTimestamps,
			"invalid deleted_at",
		)
	}

	if createdAt.Compare(updatedAt) > 0 {
		return Timestamps{}, errors.New(
			ErrInvalidTimestamps,
			"created_at is after updated_at",
			errors.WithMetadata("created_at", createdAt),
			errors.WithMetadata("updated_at", updatedAt),
		)
	}

	return Timestamps{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}, nil
}

func CreateTimestamps() Timestamps {
	now := time.Now()
	t, _ := NewTimestamps(now, now, nil)

	return t
}

func (t Timestamps) Update() Timestamps {
	t.UpdatedAt = time.Now()

	return t
}

func (t Timestamps) Delete() Timestamps {
	now := time.Now()
	t.DeletedAt = &now

	return t
}

// Serialization
type TimestampsDTO struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (t Timestamps) MarshalJSON() ([]byte, error) {
	dto := TimestampsDTO(t)

	return json.Marshal(dto)
}

func (t *Timestamps) UnmarshalJSON(data []byte) error {
	var dto TimestampsDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return err
	}

	t.CreatedAt = dto.CreatedAt
	t.UpdatedAt = dto.UpdatedAt
	t.DeletedAt = dto.DeletedAt

	return nil
}
