package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

type testItem struct {
	Name string
}

// fromNamedRow turns a row into a testItem using the first column. Rows with an empty
// first column are skipped, and rows with no columns at all are an error.
func fromNamedRow(_ testHeaderLabels, dataRow []string) (*testItem, error) {
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

type testHeaderLabels []string

func (h testHeaderLabels) Labels() []string { return h }

func TestFromSheet(t *testing.T) {
	header := testHeaderLabels([]string{"name"})

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
	header := testHeaderLabels([]string{"name"})
	rows := [][]string{{"one"}, {}}

	actual, err := FromSheet(header, rows, fromNamedRow)

	require.Error(t, err)
	// the index of the failing row should be part of the error
	assert.Contains(t, err.Error(), "row 1")
	assert.Nil(t, actual)
}

// parseNamedHeader treats a row whose first cell is "name" as a header, and a row whose
// first cell is "bad name" as a malformed header.
func parseNamedHeader(row []string) (testHeaderLabels, bool, error) {
	if len(row) == 0 {
		return nil, false, nil
	}
	switch row[0] {
	case "name":
		return row, true, nil
	case "bad name":
		return nil, true, fmt.Errorf("malformed header")
	default:
		return nil, false, nil
	}
}

// newTestFile returns an in-memory Excel file with one sheet per element of sheets, in order.
func newTestFile(t *testing.T, sheets ...[][]string) *excelize.File {
	t.Helper()
	file := excelize.NewFile()
	t.Cleanup(func() {
		assert.NoError(t, file.Close())
	})
	for i, rows := range sheets {
		sheet := fmt.Sprintf("Sheet%d", i+1)
		if i > 0 {
			_, err := file.NewSheet(sheet)
			require.NoError(t, err)
		}
		for r, row := range rows {
			cell, err := excelize.CoordinatesToCellName(1, r+1)
			require.NoError(t, err)
			require.NoError(t, file.SetSheetRow(sheet, cell, &row))
		}
	}
	return file
}

func TestFromFile(t *testing.T) {
	for name, testCase := range map[string]struct {
		sheets         [][][]string
		expectedHeader []string
		expected       []testItem
	}{
		"one sheet": {
			sheets:         [][][]string{{{"name", "extra"}, {"one"}, {"two"}}},
			expectedHeader: []string{"name", "extra"},
			expected:       []testItem{{Name: "one"}, {Name: "two"}},
		},
		"later sheet without header uses previous header": {
			sheets:         [][][]string{{{"name"}, {"one"}}, {{"two"}}},
			expectedHeader: []string{"name"},
			expected:       []testItem{{Name: "one"}, {Name: "two"}},
		},
		"later sheet with its own header": {
			sheets:         [][][]string{{{"name"}, {"one"}}, {{"name", "extra"}, {"two"}}},
			expectedHeader: []string{"name", "extra"},
			expected:       []testItem{{Name: "one"}, {Name: "two"}},
		},
		"empty sheet is ignored": {
			sheets:         [][][]string{{}, {{"name"}, {"one"}}},
			expectedHeader: []string{"name"},
			expected:       []testItem{{Name: "one"}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			file := newTestFile(t, testCase.sheets...)

			header, actual, err := FromFile(file, parseNamedHeader, fromNamedRow)
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedHeader, header)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestFromFile_Errors(t *testing.T) {
	for name, sheets := range map[string][][][]string{
		"no header":                {{{"one"}}},
		"malformed header":         {{{"bad name"}, {"one"}}},
		"malformed later header":   {{{"name"}, {"one"}}, {{"bad name"}, {"two"}}},
		"row error after a header": {{{"name"}, {"one"}, {}, {"two"}}},
	} {
		t.Run(name, func(t *testing.T) {
			file := newTestFile(t, sheets...)

			header, actual, err := FromFile(file, parseNamedHeader, fromNamedRow)
			assert.Error(t, err)
			assert.Nil(t, header)
			assert.Nil(t, actual)
		})
	}
}
