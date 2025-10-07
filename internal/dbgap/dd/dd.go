package dd

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/utils"
	"github.com/xuri/excelize/v2"
)

var logger = logging.PackageLogger("dd")

var HeaderStyle = &excelize.Style{Font: &excelize.Font{Bold: true}}

type EncodedValue struct {
	Value   string
	Meaning string
}

func NewEncodedValue(value, meaning string) EncodedValue {
	return EncodedValue{
		Value:   value,
		Meaning: meaning,
	}
}

func (v EncodedValue) String() string {
	return fmt.Sprintf("%s=%s", v.Value, v.Meaning)
}

type Column string

const VarNameColumn = Column("VARNAME")

const VarDescColumn = Column("VARDESC")

const TypeColumn = Column("TYPE")

const ValuesColumn = Column("VALUES")

const UniqueKeyColumn = Column("UNIQUEKEY")

type Spec struct {
	FileName  string
	SheetName string
	Header    []Column
	Rows      [][]any
}

func populateRow[T any](f *excelize.File, sheetName string, rowNumber int, row []T, colWidths map[int]int) (map[int]int, error) {
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

func Populate(f *excelize.File, sheet string, spec Spec) error {

	// Write Header
	colWidths, err := populateRow(f, sheet, 1, spec.Header, nil)
	if err != nil {
		return err
	}

	// Write rows
	for r, row := range spec.Rows {
		colWidths, err = populateRow(f, sheet, r+2, row, colWidths)
		if err != nil {
			return err
		}
	}

	// Style header bold
	if style, err := f.NewStyle(HeaderStyle); err != nil {
		return fmt.Errorf("error adding header style to DD file: %w", err)
	} else {
		if err := f.SetRowStyle(sheet, 1, 1, style); err != nil {
			return fmt.Errorf("error setting header style to DD file: %w", err)
		}
	}

	// Apply column widths (+2 padding)
	for c, w := range colWidths {
		if colName, err := excelize.ColumnNumberToName(c + 1); err != nil {
			return fmt.Errorf("error getting column name of DD file: %w", err)
		} else {
			if err := f.SetColWidth(sheet, colName, colName, float64(w+2)); err != nil {
				return fmt.Errorf("error setting width of column %s in DD file: %w", colName, err)
			}
		}
	}
	return nil
}

func Write(path string, spec Spec) error {
	subjectConsentDDFile := excelize.NewFile()
	defer utils.CloseExcelFile(subjectConsentDDFile, logger)

	if err := subjectConsentDDFile.SetSheetName("Sheet1", spec.SheetName); err != nil {
		return fmt.Errorf("error setting %s sheet name: %w", spec.FileName, err)
	}
	if err := Populate(subjectConsentDDFile, spec.SheetName, spec); err != nil {
		return fmt.Errorf("error populating %s: %w", spec.FileName, err)
	}
	if err := subjectConsentDDFile.SaveAs(path); err != nil {
		return fmt.Errorf("error writing %s to %s: %w", spec.FileName, path, err)
	}
	return nil
}
