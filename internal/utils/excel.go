package utils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log/slog"
)

func CloseExcelFile(inputFile *excelize.File, logger *slog.Logger) {
	if err := inputFile.Close(); err != nil {
		logger.Warn("error closing Excel file",
			slog.String("path", inputFile.Path),
			slog.Any("error", err))
	}
}

type IsHeaderRowFunc func(row []string) bool

type FromRowFunc[T any] func(header []string, dataRow []string) (T, error)

func FromSheet[T any](header []string, rows [][]string, fromRow FromRowFunc[T]) ([]T, error) {
	items := make([]T, 0, len(rows))

	for _, row := range rows {
		if sample, err := fromRow(header, row); err != nil {
			return nil, err
		} else {
			items = append(items, sample)
		}
	}
	return items, nil
}

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
