package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testItem struct {
	Name string
}

// fromNamedRow turns a row into a testItem using the first column. Rows with an empty
// first column are skipped, and rows with no columns at all are an error.
func fromNamedRow(_ []string, dataRow []string) (*testItem, error) {
	if len(dataRow) == 0 {
		return nil, fmt.Errorf("row is empty")
	}
	if len(dataRow[0]) == 0 {
		return nil, nil
	}
	return &testItem{Name: dataRow[0]}, nil
}

func TestIsBlankRow(t *testing.T) {
	for name, testCase := range map[string]struct {
		row      []string
		expected bool
	}{
		"no cells":          {row: []string{}, expected: true},
		"nil row":           {row: nil, expected: true},
		"empty cells":       {row: []string{"", ""}, expected: true},
		"whitespace cells":  {row: []string{" ", "\t", "\n"}, expected: true},
		"one non-empty":     {row: []string{"", "value", ""}, expected: false},
		"all non-empty":     {row: []string{"one", "two"}, expected: false},
		"zero looks filled": {row: []string{"0"}, expected: false},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, IsBlankRow(testCase.row))
		})
	}
}

func TestFromSheet(t *testing.T) {
	header := []string{"name"}

	for name, testCase := range map[string]struct {
		rows     [][]string
		expected []testItem
	}{
		"no rows": {
			rows:     [][]string{},
			expected: []testItem{},
		},
		"no skipped rows": {
			rows:     [][]string{{"one"}, {"two"}},
			expected: []testItem{{Name: "one"}, {Name: "two"}},
		},
		"skipped rows are left out": {
			rows:     [][]string{{"one"}, {""}, {"two"}, {""}},
			expected: []testItem{{Name: "one"}, {Name: "two"}},
		},
		"every row skipped": {
			rows:     [][]string{{""}, {""}},
			expected: []testItem{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual, err := FromSheet(header, testCase.rows, fromNamedRow)
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestFromSheet_Error(t *testing.T) {
	header := []string{"name"}
	rows := [][]string{{"one"}, {}}

	actual, err := FromSheet(header, rows, fromNamedRow)

	require.Error(t, err)
	// the index of the failing row should be part of the error
	assert.Contains(t, err.Error(), "row 1")
	assert.Nil(t, actual)
}
