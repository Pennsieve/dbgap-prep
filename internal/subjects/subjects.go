package subjects

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/utils"
	"github.com/xuri/excelize/v2"
)

var logger = logging.PackageLogger("subjects")

const FileName = "subjects.xlsx"

const SourceSubjectIDdbGapColumn = "SOURCE_SUBJECT_ID_dbGaP"

const IDIndex = 0
const IDLabel = "subject id"

// ValidIDPrefix is the prefix that distinguishes a row in the subjects file that
// should be submitted to dbGaP from one that should not.
// Matching against it ignores case and surrounding whitespace, but
// IDs are stored as they appear in the file so that they still match the samples file.
const ValidIDPrefix = "sub-"
const SexLabel = "sex"

type Subject struct {
	ID                   string
	Sex                  string
	SourceSubjectIDdbGap string
	// Values maps header labels to corresponding values for this row
	Values map[string]string
}

func (s Subject) GetValue(key string) (string, bool) {
	value, ok := s.Values[key]
	return value, ok
}

func (s Subject) SearchValue(key string) (string, bool) {
	for label, value := range s.Values {
		if strings.EqualFold(key, label) {
			return value, true
		}
	}
	return "", false
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

// HasValidIDPrefix returns true if the given subject ID starts with ValidIDPrefix,
// ignoring case and any surrounding whitespace.
func HasValidIDPrefix(id string) bool {
	trimmed := strings.TrimSpace(id)
	if len(trimmed) < len(ValidIDPrefix) {
		return false
	}
	return strings.EqualFold(trimmed[:len(ValidIDPrefix)], ValidIDPrefix)
}

// FromRow converts the given non-header row to a Subject.
// Blank rows and rows whose ID does not have the ValidIDPrefix prefix are skipped:
// for those, a nil Subject and a nil error are returned so that the caller skips the row
// without failing.
func FromRow(header []string, row []string) (*Subject, error) {
	if IsHeaderRow(row) {
		return nil, fmt.Errorf("subjects row is a header")
	}
	if utils.IsBlankRow(row) {
		logger.Info("skipping blank subjects row")
		return nil, nil
	}
	if len(row) < IDIndex+1 {
		return nil, fmt.Errorf("subjects row is too short to contain required columns")
	}
	if !HasValidIDPrefix(row[IDIndex]) {
		logger.Info("skipping row; subject ID does not have the expected prefix",
			slog.String("id", row[IDIndex]),
			slog.String("expectedPrefix", ValidIDPrefix))
		return nil, nil
	}
	// the ID column is always kept out of values
	values := make(map[string]string, len(row)-1)
	subject := Subject{
		ID:     row[IDIndex],
		Values: values,
	}

	for i, label := range header {
		if i == IDIndex {
			//skip this since it is already part of the struct
		} else if i < len(row) {
			// excelize does not give us empty cells beyond the last non-empty cell
			if strings.EqualFold(label, SourceSubjectIDdbGapColumn) {
				subject.SourceSubjectIDdbGap = row[i]
			} else if strings.EqualFold(label, SexLabel) {
				subject.Sex = row[i]
			} else {
				values[label] = row[i]
			}
		}
	}
	logger.Info("found subject", subject.LogGroup())
	return &subject, nil
}

func FromFile(subjectsFile *excelize.File) ([]string, []Subject, error) {
	return utils.FromFile(subjectsFile, IsHeaderRow, FromRow)
}
