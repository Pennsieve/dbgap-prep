package utils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log/slog"
)

var HeaderStyle = &excelize.Style{Font: &excelize.Font{Bold: true}}

func CloseExcelFile(inputFile *excelize.File, logger *slog.Logger) {
	if err := inputFile.Close(); err != nil {
		logger.Warn("error closing Excel file",
			slog.String("path", inputFile.Path),
			slog.Any("error", err))
	}
}

// Reading

type IsHeaderRowFunc func(row []string) bool

// FromRowFunc is a function that, given a header and dataRow, turn that dataRow into a T.
type FromRowFunc[T any] func(header []string, dataRow []string) (T, error)

// FromSheet returns a slice of T, created from header and rows using fromRow.
func FromSheet[T any](header []string, rows [][]string, fromRow FromRowFunc[T]) ([]T, error) {
	items := make([]T, 0, len(rows))

	for i, row := range rows {
		sample, err := fromRow(header, row)
		if err != nil {
			return nil, fmt.Errorf("error converting row %d: %w", i, err)
		}
		items = append(items, sample)

	}
	return items, nil
}

// FromFile returns a slice of T, using isHeaderRow to determine the header row, and fromRow to turn non-header rows into Ts.
// The first returned argument are the header labels.
func FromFile[T any](file *excelize.File, isHeaderRow IsHeaderRowFunc, fromRow FromRowFunc[T]) ([]string, []T, error) {
	var allItems []T
	var header []string

	for _, sheetName := range file.GetSheetList() {
		rows, err := file.GetRows(sheetName)
		if err != nil {
			return nil, nil, fmt.Errorf("error getting rows from %s, sheet %s: %w", file.Path, sheetName, err)
		}

		if len(rows) > 0 {
			dataRows := rows
			maybeHeader := rows[0]
			if isHeaderRow(maybeHeader) {
				header = maybeHeader
				dataRows = rows[1:]
			} else if header == nil {
				// First row in sheet is not a header, and there is no header from
				// previous sheets, so return error
				return nil, nil, fmt.Errorf("no header found for sheet %s", sheetName)
			}
			items, err := FromSheet(header, dataRows, fromRow)
			if err != nil {
				return nil, nil, fmt.Errorf("error getting items from %s, sheet %s: %w", file.Path, sheetName, err)
			}
			allItems = append(allItems, items...)
		}
	}

	return header, allItems, nil
}

// Writing

// PopulateRow writes the slice of T to the given sheet and file at rowNumber (1-based).
// The returned map maps column numbers (0-based) to widths. If the passed colWidths is not-nil, the returned map is an updated
// version of colWidths.
func PopulateRow[T any](f *excelize.File, sheetName string, rowNumber int, row []T, colWidths map[int]int) (map[int]int, error) {
	if colWidths == nil {
		colWidths = map[int]int{}
	}

	// get starting cell for row
	cell, err := excelize.CoordinatesToCellName(1, rowNumber)
	if err != nil {
		return nil, fmt.Errorf("error getting DD file cell name: %w", err)
	}
	if err := f.SetSheetRow(sheetName, cell, &row); err != nil {
		return nil, fmt.Errorf("error setting DD file row %d: %w", rowNumber, err)
	}

	// update column widths
	for c, v := range row {
		str := fmt.Sprint(v) // convert to string for length
		if len(str) > colWidths[c] {
			colWidths[c] = len(str)
		}
	}

	return colWidths, nil
}
