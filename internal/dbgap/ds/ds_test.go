package ds

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToRows(t *testing.T) {
	type TestItem struct {
		ID string
	}

	toRow := func(_ []string, item TestItem) []string {
		return []string{item.ID}
	}

	var items []TestItem
	for i := 0; i < 13; i++ {
		items = append(items, TestItem{ID: uuid.NewString()})
	}

	rows := ToRows(nil, items, toRow)

	require.Len(t, rows, len(items))

	for i, row := range rows {
		require.Len(t, row, 1)
		assert.Equal(t, items[i].ID, row[0])
	}
}
