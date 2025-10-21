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

	if err := colWidths.SetWidths(f, sheet); err != nil {
		return fmt.Errorf("error setting column widths of DD file: %w", err)
	}

	return nil
}
