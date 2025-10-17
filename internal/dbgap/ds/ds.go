package ds

import (
	"encoding/csv"
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/utils"
	"github.com/xuri/excelize/v2"
	"log/slog"
	"os"
)

var logger = logging.PackageLogger("ds")

type Spec struct {
	FileName  string
	Variables []dd.Variable
}

func Write(path string, spec Spec, rows [][]string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating DS file %s: %w", path, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			logger.Warn("error closing DS file", slog.String("path", path),
				slog.Any("error", err))
		}
	}()

	tsvWriter := csv.NewWriter(file)
	tsvWriter.Comma = '\t'
	// some users may use windows
	tsvWriter.UseCRLF = true

	records := make([][]string, 0, len(rows)+1)
	header := dd.VariableNames(spec.Variables)
	records = append(records, header)
	records = append(records, rows...)
	if err := tsvWriter.WriteAll(records); err != nil {
		return fmt.Errorf("error writing records to DS file %s: %s", path, err)
	}

	return nil
}

func WriteXLSX(path string, spec Spec, rows [][]string) error {
	f := excelize.NewFile()
	defer utils.CloseExcelFile(f, logger)

	sheet := "Sheet1"
	header := dd.VariableNames(spec.Variables)

	// Write Header
	colWidths, err := utils.PopulateRow(f, sheet, 1, header, nil)
	if err != nil {
		return err
	}

	// Write rows
	for r, row := range rows {
		colWidths, err = utils.PopulateRow(f, sheet, r+2, row, colWidths)
		if err != nil {
			return err
		}
	}

	// Style header bold
	if style, err := f.NewStyle(utils.HeaderStyle); err != nil {
		return fmt.Errorf("error adding header style to DS file: %w", err)
	} else {
		if err := f.SetRowStyle(sheet, 1, 1, style); err != nil {
			return fmt.Errorf("error setting header style to DS file: %w", err)
		}
	}

	// Apply column widths (+2 padding)
	for c, w := range colWidths {
		if colName, err := excelize.ColumnNumberToName(c + 1); err != nil {
			return fmt.Errorf("error getting column name of DS file: %w", err)
		} else {
			if err := f.SetColWidth(sheet, colName, colName, float64(w+2)); err != nil {
				return fmt.Errorf("error setting width of column %s in DS file: %w", colName, err)
			}
		}
	}

	if err := f.SaveAs(path); err != nil {
		return fmt.Errorf("error saving XLSX file to %s: %w", path, err)
	}
	return nil
}

type ToRowFunc[T any] func(variableNames []string, item T) []string

func ToRows[T any](variables []dd.Variable, items []T, toRow ToRowFunc[T]) [][]string {
	variableNames := dd.VariableNames(variables)

	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, toRow(variableNames, item))
	}
	return rows
}
