package subjects

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/utils"
	"github.com/xuri/excelize/v2"
)

var logger = logging.PackageLogger("subjects")

const FileName = "subjects.xlsx"

const SourceSubjectIDdbGapLabel = "SOURCE_SUBJECT_ID_dbGaP"

const IDIndex = 0
const IDLabel = "subject id"

// ValidIDPrefix is the prefix that distinguishes a row in the subjects file that
// should be submitted to dbGaP from one that should not.
// Matching against it ignores case and surrounding whitespace, but
// IDs are stored as they appear in the file so that they still match the samples file.
const ValidIDPrefix = "sub-"
const SexLabel = "sex"

const minLength = IDIndex + 1

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

type Header struct {
	labels []string
	// sexIndex is -1 if Header does not contain a sex label
	sexIndex int
	// sourceSubjectIDdbGapIndex is -1 if Header does not contain a SOURCE_SUBJECT_ID_dbGaP label
	sourceSubjectIDdbGapIndex int
}

func (h Header) Labels() []string {
	return h.labels
}

func ParseHeaderRow(row []string) (Header, bool, error) {
	if len(row) == 0 || !utils.EqualIgnoreCaseTrimmed(row[IDIndex], IDLabel) {
		return Header{}, false, nil
	}

	sexIndex, sourceSubjectIDdbGapIndex := -1, -1
	for i, label := range row {
		if utils.EqualIgnoreCaseTrimmed(label, SexLabel) {
			sexIndex = i
		} else if utils.EqualIgnoreCaseTrimmed(label, SourceSubjectIDdbGapLabel) {
			sourceSubjectIDdbGapIndex = i
		}
	}

	// defensive copy that uses the canonical labels even if actual header had different case
	labels := slices.Clone(row)
	labels[IDIndex] = IDLabel
	if sexIndex >= 0 {
		labels[sexIndex] = SexLabel
	}
	if sourceSubjectIDdbGapIndex >= 0 {
		labels[sourceSubjectIDdbGapIndex] = SourceSubjectIDdbGapLabel
	}

	return Header{labels: labels, sexIndex: sexIndex, sourceSubjectIDdbGapIndex: sourceSubjectIDdbGapIndex}, true, nil
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
func FromRow(header Header, row []string) (*Subject, error) {
	if utils.IsBlankRow(row) {
		logger.Info("skipping blank subjects row")
		return nil, nil
	}
	if len(row) < minLength {
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

	for i, label := range header.Labels() {
		if i >= len(row) {
			break // excelize does not give us empty cells beyond the last non-empty cell
		}
		switch i {
		case IDIndex:
			// skip this since it is already part of the struct
		case header.sexIndex:
			subject.Sex = row[i]
		case header.sourceSubjectIDdbGapIndex:
			subject.SourceSubjectIDdbGap = row[i]
		default:
			values[label] = row[i]
		}
	}
	logger.Info("found subject", subject.LogGroup())
	return &subject, nil
}

func FromFile(subjectsFile *excelize.File) ([]string, []Subject, error) {
	return utils.FromFile(subjectsFile, ParseHeaderRow, FromRow)
}
