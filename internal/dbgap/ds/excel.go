package ds

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/utils"
	"github.com/xuri/excelize/v2"
	"path/filepath"
)

type XLSXWriter struct {
	path string
}

func NewXLSXWriter(outputDirectory string, baseFileName string) *XLSXWriter {
	return &XLSXWriter{path: filepath.Join(outputDirectory, fmt.Sprintf("%s.xlsx", baseFileName))}
}

func (x *XLSXWriter) Path() string {
	return x.path
}

func (x *XLSXWriter) Write(spec Spec, rows [][]string) error {
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

	if err := f.SaveAs(x.path); err != nil {
		return fmt.Errorf("error saving XLSX file to %s: %w", x.path, err)
	}
	return nil
}
