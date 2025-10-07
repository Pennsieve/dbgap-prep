package samples

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/xuri/excelize/v2"
	"log/slog"
)

var logger = logging.PackageLogger("samples")

const FileName = "samples.xlsx"

const IDIndex = 0
const IDLabel = "sample id"
const SubjectIDIndex = 1
const SubjectIDLabel = "subject id"

type Sample struct {
	ID        string
	SubjectID string
	// Values maps header labels to corresponding values for this row
	Values map[string]string
}

func (s Sample) String() string {
	return fmt.Sprintf("sample: id = [%s], subject id = [%s], valueCount = %d",
		s.ID,
		s.SubjectID,
		len(s.Values))
}

func (s Sample) LogGroup() slog.Attr {
	return slog.Group("sample",
		slog.String("id", s.ID),
		slog.String("subjectId", s.SubjectID),
		slog.Int("valueCount", len(s.Values)),
	)
}

func (s Sample) HasSubject() bool {
	return len(s.SubjectID) > 0
}

func IsHeaderRow(row []string) bool {
	return len(row) > 0 && row[0] == IDLabel
}

// FromRow converts the given non-empty, non-header row to a Sample
func FromRow(header []string, row []string) (Sample, error) {
	if IsHeaderRow(row) {
		return Sample{}, fmt.Errorf("samples row is a header")
	}
	if len(row) < 2 {
		return Sample{}, fmt.Errorf("samples row is too short to contain sample and subject ids")
	}
	values := make(map[string]string, len(row)-2)
	sample := Sample{
		ID:        row[IDIndex],
		SubjectID: row[SubjectIDIndex],
		Values:    values,
	}

	for i, label := range header {
		if i == IDIndex || i == SubjectIDIndex {
			//skip these since they are already part of the struct
		} else if i < len(row) {
			// excelize does not give us empty cells beyond the last non-empty cell
			values[label] = row[i]
		} else {
			// maybe we'll have to distinguish between a missing value for a real label
			// and a bad label?
			values[label] = ""
		}
	}
	logger.Info("found sample", sample.LogGroup())
	return sample, nil
}

func FromSheet(header []string, rows [][]string) ([]Sample, error) {
	samps := make([]Sample, 0, len(rows))

	for _, row := range rows {
		if sample, err := FromRow(header, row); err != nil {
			return nil, err
		} else {
			samps = append(samps, sample)
		}
	}
	return samps, nil
}

func FromFile(file *excelize.File) ([]string, []Sample, error) {
	var allSamps []Sample
	var header []string
	for _, sheetName := range file.GetSheetList() {
		rows, err := file.GetRows(sheetName)
		if err != nil {
			return nil, nil, fmt.Errorf("error getting rows from %s, sheet %s: %w", file.Path, sheetName, err)
		}

		if len(rows) > 0 {
			dataRows := rows
			maybeHeader := rows[0]
			if IsHeaderRow(maybeHeader) {
				header = maybeHeader
				dataRows = rows[1:]
			} else if header == nil {
				// First row in sheet is not a header, and there is no header from
				// previous sheets, so return error
				return nil, nil, fmt.Errorf("no header found for sheet %s", sheetName)
			}
			samps, err := FromSheet(header, dataRows)
			if err != nil {
				return nil, nil, fmt.Errorf("error getting samples from %s, sheet %s: %w", file.Path, sheetName, err)
			}
			allSamps = append(allSamps, samps...)
		}
	}
	return header, allSamps, nil
}
