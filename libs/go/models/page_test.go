package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPage(t *testing.T) {
	t.Run("page", func(t *testing.T) {
		items := []string{"item1", "item2", "item3"}
		total := 5

		page := NewPage(items, total)

		assert.Equal(t, items, page.Items)
		assert.Equal(t, len(items), page.Count)
		assert.Equal(t, 5, page.Total)

		cursor := OffsetLimitCursor{
			Offset: 0,
			Limit:  10,
		}
		page = page.WithCursor(cursor)

		assert.True(t, page.HasNext())
		assert.IsType(t, "", *page.Cursor)
	})

	t.Run("cursor", func(t *testing.T) {
		t.Run("offset limit cursor", func(t *testing.T) {
			cursor := OffsetLimitCursor{
				Offset: 80,
				Limit:  75,
			}

			encoded := cursor.Encode()
			assert.NotEmpty(t, encoded)

			decodedCursor, err := DecodeOffsetLimitCursor(encoded)
			assert.NoError(t, err)
			assert.Equal(t, cursor, decodedCursor)
		})

		t.Run("page cursor", func(t *testing.T) {
			cursor := PageCursor{
				Page: 3,
			}

			encoded := cursor.Encode()
			assert.NotEmpty(t, encoded)

			decodedCursor, err := DecodePageCursor(encoded)
			assert.NoError(t, err)
			assert.Equal(t, cursor, decodedCursor)
		})

		t.Run("key cursor", func(t *testing.T) {
			cursor := KeyCursor{
				Key: "key",
			}

			encoded := cursor.Encode()
			assert.NotEmpty(t, encoded)

			decodedCursor, err := DecodeKeyCursor(encoded)
			assert.NoError(t, err)
			assert.Equal(t, cursor, decodedCursor)
		})
	})
}
