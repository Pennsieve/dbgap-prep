package samples

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/utils"
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
func FromRow(header []string, row []string) (Sample, []string, error) {
	if IsHeaderRow(row) {
		return Sample{}, nil, fmt.Errorf("samples row is a header")
	}
	if len(row) < 2 {
		return Sample{}, nil, fmt.Errorf("samples row is too short to contain sample and subject ids")
	}
	values := make(map[string]string, len(row)-2)
	nonEmptyKeys := make([]string, 0, len(header))

	sample := Sample{
		ID:        row[IDIndex],
		SubjectID: row[SubjectIDIndex],
		Values:    values,
	}

	for i, label := range header {
		if i == IDIndex || i == SubjectIDIndex {
			//skip these since they are already part of the struct.
			// but always include them in the nonempty keys
			nonEmptyKeys = append(nonEmptyKeys, label)
		} else if i < len(row) {
			// excelize does not give us empty cells beyond the last non-empty cell
			value := row[i]
			if len(value) > 0 {
				nonEmptyKeys = append(nonEmptyKeys, label)
			}
			values[label] = value
		}
	}
	logger.Info("found sample", sample.LogGroup())
	return sample, nonEmptyKeys, nil
}

func FromFile(file *excelize.File) ([]string, []Sample, error) {
	return utils.FromFile(file, IsHeaderRow, FromRow)
}
