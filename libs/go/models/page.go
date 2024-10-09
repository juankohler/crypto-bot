package models

import (
	"encoding/base64"
	"net/url"
	"strconv"

	"github.com/juankohler/crypto-bot/libs/go/errors"
)

var (
	ErrPageInvalidCursor = errors.Define("page.invalid_cursor")
)

// Page
type Page[T any] struct {
	Items  []T     `json:"items"`
	Count  int     `json:"count"`
	Total  int     `json:"total"`
	Cursor *string `json:"cursor"`
}

func NewPage[T any](items []T, total int) Page[T] {
	return Page[T]{
		Items: items,
		Count: len(items),
		Total: total,
	}
}

func (p Page[T]) HasNext() bool {
	return p.Cursor != nil
}

func (p Page[T]) WithCursor(cursor Cursor) Page[T] {
	encoded := cursor.Encode()
	encodedCursor := &encoded

	return Page[T]{
		Items:  p.Items,
		Count:  p.Count,
		Total:  p.Total,
		Cursor: encodedCursor,
	}
}

// Cursors
type Cursor interface {
	Encode() string
}

type OffsetLimitCursor struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (c OffsetLimitCursor) Encode() string {
	v := url.Values{}

	v.Set("offset", strconv.Itoa(c.Offset))
	v.Set("limit", strconv.Itoa(c.Limit))

	return base64.URLEncoding.EncodeToString([]byte(v.Encode()))
}

func DecodeOffsetLimitCursor(cursor string) (OffsetLimitCursor, error) {
	data, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return OffsetLimitCursor{}, errors.Wrap(
			ErrPageInvalidCursor,
			err,
			"invalid cursor",
			errors.WithMetadata("cursor", cursor),
		)
	}

	v, err := url.ParseQuery(string(data))
	if err != nil {
		return OffsetLimitCursor{}, errors.Wrap(
			ErrPageInvalidCursor,
			err,
			"invalid cursor",
			errors.WithMetadata("cursor", cursor),
		)
	}

	offset, err := strconv.Atoi(v.Get("offset"))
	if err != nil {
		return OffsetLimitCursor{}, errors.Wrap(
			ErrPageInvalidCursor,
			err,
			"invalid cursor",
			errors.WithMetadata("cursor", cursor),
		)
	}

	limit, err := strconv.Atoi(v.Get("limit"))
	if err != nil {
		return OffsetLimitCursor{}, errors.Wrap(
			ErrPageInvalidCursor,
			err,
			"invalid cursor",
			errors.WithMetadata("cursor", cursor),
		)
	}

	return OffsetLimitCursor{
		Offset: offset,
		Limit:  limit,
	}, nil
}

type PageCursor struct {
	Page int `json:"page"`
}

func (c PageCursor) Encode() string {
	v := url.Values{}

	v.Set("page", strconv.Itoa(c.Page))

	return base64.URLEncoding.EncodeToString([]byte(v.Encode()))
}

func DecodePageCursor(cursor string) (PageCursor, error) {
	data, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return PageCursor{}, errors.Wrap(
			ErrPageInvalidCursor,
			err,
			"invalid cursor",
			errors.WithMetadata("cursor", cursor),
		)
	}

	v, err := url.ParseQuery(string(data))
	if err != nil {
		return PageCursor{}, errors.Wrap(
			ErrPageInvalidCursor,
			err,
			"invalid cursor",
			errors.WithMetadata("cursor", cursor),
		)
	}

	page, err := strconv.Atoi(v.Get("page"))
	if err != nil {
		return PageCursor{}, errors.Wrap(
			ErrPageInvalidCursor,
			err,
			"invalid cursor",
			errors.WithMetadata("cursor", cursor),
		)
	}

	return PageCursor{
		Page: page,
	}, nil
}

type KeyCursor struct {
	Key string `json:"key"`
}

func (c KeyCursor) Encode() string {
	v := url.Values{}

	v.Set("key", c.Key)

	return base64.URLEncoding.EncodeToString([]byte(v.Encode()))
}

func DecodeKeyCursor(cursor string) (KeyCursor, error) {
	data, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return KeyCursor{}, errors.Wrap(
			ErrPageInvalidCursor,
			err,
			"invalid cursor",
			errors.WithMetadata("cursor", cursor),
		)
	}

	v, err := url.ParseQuery(string(data))
	if err != nil {
		return KeyCursor{}, errors.Wrap(
			ErrPageInvalidCursor,
			err,
			"invalid cursor",
			errors.WithMetadata("cursor", cursor),
		)
	}

	key := v.Get("key")

	return KeyCursor{
		Key: key,
	}, nil
}
