package subjects

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/utils"
	"github.com/xuri/excelize/v2"
	"log/slog"
)

var logger = logging.PackageLogger("subjects")

const FileName = "subjects.xlsx"

const IDIndex = 0
const IDLabel = "subject id"
const SexIndex = 4
const SexLabel = "sex"

type Subject struct {
	ID  string
	Sex string
	// Values maps header labels to corresponding values for this row
	Values map[string]string
}

func (s Subject) String() string {
	return fmt.Sprintf("subject: id = [%s], sex = [%s], valueCount = %d",
		s.ID,
		s.Sex,
		len(s.Values),
	)
}

func (s Subject) LogGroup() slog.Attr {
	return slog.Group("subject",
		slog.String("id", s.ID),
		slog.String("sex", s.Sex),
		slog.Int("valueCount", len(s.Values)),
	)
}

func IsHeaderRow(row []string) bool {
	return len(row) > 0 &&
		row[IDIndex] == IDLabel
}

// FromRow converts the given non-empty, non-header row to a Subject
func FromRow(header []string, row []string) (Subject, []string, error) {
	if IsHeaderRow(row) {
		return Subject{}, nil, fmt.Errorf("subjects row is a header")
	}
	if len(row) < IDIndex+1 {
		return Subject{}, nil, fmt.Errorf("subjects row is too short to contain required columns")
	}
	var sex string
	if len(row) > SexIndex {
		sex = row[SexIndex]
	}
	values := make(map[string]string, len(row)-2)
	nonEmptyKeys := make([]string, len(header))
	subject := Subject{
		ID:     row[IDIndex],
		Sex:    sex,
		Values: values,
	}

	for i, label := range header {
		if i == IDIndex || i == SexIndex {
			//skip these since they are already part of the struct
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
	logger.Info("found subject", subject.LogGroup())
	return subject, nonEmptyKeys, nil
}

func FromFile(subjectsFile *excelize.File) ([]string, []Subject, error) {
	return utils.FromFile(subjectsFile, IsHeaderRow, FromRow)
}
