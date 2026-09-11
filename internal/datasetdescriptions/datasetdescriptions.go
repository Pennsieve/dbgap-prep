package datasetdescriptions

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/xuri/excelize/v2"
)

var logger = logging.PackageLogger("datasetdescriptions")

const FileName = "dataset_description.xlsx"
const TitleCellAddress = "D8"
const IdentifierDescriptionRowIndex = 32
const IdentifierRowIndex = 34

// DOIURLIdentifierDescriptions are the descriptions that exist in practice.
// Deliberately all lower case here for the comparison in the code.
var DOIURLIdentifierDescriptions = map[string]bool{
	"doi for the first version of this dataset": true,
	"doi for this dataset":                      true,
	"doi for first version of this dataset":     true,
}

type DatasetDescription struct {
	DOIURL string
	Title  string
}

func (dd DatasetDescription) LogGroup() slog.Attr {
	return slog.Group("datasetDescription",
		slog.String("title", dd.Title),
		slog.String("doiUrl", dd.DOIURL),
	)
}

func FromFile(file *excelize.File) (DatasetDescription, error) {
	datasetDescription := DatasetDescription{}
	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return datasetDescription, fmt.Errorf("no sheets in %s", FileName)
	}
	sheet := sheets[0]

	// Get Title
	title, err := file.GetCellValue(sheet, TitleCellAddress)
	if err != nil {
		return datasetDescription, fmt.Errorf("error getting title cell %s from %s: %w",
			TitleCellAddress,
			FileName,
			err)
	}
	datasetDescription.Title = title

	//Get DOI URL
	rows, err := file.GetRows(sheet)
	if err != nil {
		return datasetDescription, fmt.Errorf("error getting rows from %s: %w", FileName, err)
	}

	// relying on IdentifierRowIndex > IdentifierDescriptionRowIndex
	if len(rows) < IdentifierRowIndex+1 {
		return datasetDescription, fmt.Errorf("not enough rows for in %s for identifier description", FileName)
	}
	identifierDescriptionRow := rows[IdentifierDescriptionRowIndex]
	doiURLColumnIndex := -1
	for i := range identifierDescriptionRow {
		if DOIURLIdentifierDescriptions[strings.ToLower(identifierDescriptionRow[i])] {
			doiURLColumnIndex = i
			break
		}
	}
	if doiURLColumnIndex == -1 {
		return datasetDescription, fmt.Errorf("DOI URL identifier description not found in %s", FileName)
	}
	identifierRow := rows[IdentifierRowIndex]
	if len(identifierRow) < doiURLColumnIndex+1 {
		return datasetDescription, fmt.Errorf("not enough columns for DOI URL in %s", FileName)
	}
	datasetDescription.DOIURL = identifierRow[doiURLColumnIndex]

	return datasetDescription, nil
}
