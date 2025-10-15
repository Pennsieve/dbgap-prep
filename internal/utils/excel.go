package utils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log/slog"
	"maps"
	"slices"
)

func CloseExcelFile(inputFile *excelize.File, logger *slog.Logger) {
	if err := inputFile.Close(); err != nil {
		logger.Warn("error closing Excel file",
			slog.String("path", inputFile.Path),
			slog.Any("error", err))
	}
}

type IsHeaderRowFunc func(row []string) bool

// FromRowFunc is a function that, given a header and dataRow, turn that dataRow into a T.
// The second return argument should be those values of header that have non-empty values in dataRow.
type FromRowFunc[T any] func(header []string, dataRow []string) (T, []string, error)

// FromSheet returns a slice of T, created from header and rows using fromRow.
// The second argument is a map of header strings, acting as the set of header strings that have at least one non-empty value in their column in the sheet.
func FromSheet[T any](header []string, rows [][]string, fromRow FromRowFunc[T]) ([]T, map[string]bool, error) {
	items := make([]T, 0, len(rows))
	nonEmptyKeys := map[string]bool{}

	for _, row := range rows {
		if sample, rowNonEmptyKeys, err := fromRow(header, row); err != nil {
			return nil, nil, err
		} else {
			for _, key := range rowNonEmptyKeys {
				nonEmptyKeys[key] = true
			}
			items = append(items, sample)
		}
	}
	return items, nonEmptyKeys, nil
}

// FromFile returns a slice of T, using isHeaderRow to determine the header row, and fromRow to turn non-header rows into Ts.
// The first returned argument are those header labels that contain at least one non-empty value in their column.
func FromFile[T any](file *excelize.File, isHeaderRow IsHeaderRowFunc, fromRow FromRowFunc[T]) ([]string, []T, error) {
	var allItems []T
	var header []string
	nonEmptyKeys := map[string]bool{}

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
			items, sheetNonEmptyKeys, err := FromSheet(header, dataRows, fromRow)
			if err != nil {
				return nil, nil, fmt.Errorf("error getting items from %s, sheet %s: %w", file.Path, sheetName, err)
			}
			maps.Copy(nonEmptyKeys, sheetNonEmptyKeys)
			allItems = append(allItems, items...)
		}
	}

	header = slices.DeleteFunc(header, func(key string) bool {
		_, nonEmpty := nonEmptyKeys[key]
		return !nonEmpty
	})
	return header, allItems, nil
}
