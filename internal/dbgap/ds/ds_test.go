package ds

import (
	"testing"

	"github.com/google/uuid"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToRows(t *testing.T) {
	type TestItem struct {
		ID string
	}

	toRow := func(_ []dd.Variable, item TestItem) []string {
		return []string{item.ID}
	}

	var items []TestItem
	for range 13 {
		items = append(items, TestItem{ID: uuid.NewString()})
	}

	rows := ToRows(nil, items, toRow)

	require.Len(t, rows, len(items))

	for i, row := range rows {
		require.Len(t, row, 1)
		assert.Equal(t, items[i].ID, row[0])
	}
}
