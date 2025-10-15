package dd

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/utils"
	"github.com/xuri/excelize/v2"
)

var logger = logging.PackageLogger("dd")

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

func Populate(f *excelize.File, sheet string, spec Spec) error {

	// Write Header
	colWidths, err := utils.PopulateRow(f, sheet, 1, spec.Header, nil)
	if err != nil {
		return err
	}

	// Write rows
	for r, row := range spec.Rows {
		colWidths, err = utils.PopulateRow(f, sheet, r+2, row, colWidths)
		if err != nil {
			return err
		}
	}

	// Style header bold
	if style, err := f.NewStyle(utils.HeaderStyle); err != nil {
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
	ddFile := excelize.NewFile()
	defer utils.CloseExcelFile(ddFile, logger)

	if err := ddFile.SetSheetName("Sheet1", spec.SheetName); err != nil {
		return fmt.Errorf("error setting %s sheet name: %w", spec.FileName, err)
	}
	if err := Populate(ddFile, spec.SheetName, spec); err != nil {
		return fmt.Errorf("error populating %s: %w", spec.FileName, err)
	}
	if err := ddFile.SaveAs(path); err != nil {
		return fmt.Errorf("error writing %s to %s: %w", spec.FileName, path, err)
	}
	return nil
}
