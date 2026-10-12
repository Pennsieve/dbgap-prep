package samples

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/utils"
	"github.com/xuri/excelize/v2"
)

var logger = logging.PackageLogger("samples")

const FileName = "samples.xlsx"

const IDIndex = 0
const IDLabel = "sample id"
const SubjectIDIndex = 1
const SubjectIDLabel = "subject id"
const DbGapStudyIDLabel = "dbGaP study ID"

const minLength = max(IDIndex, SubjectIDIndex) + 1

type Sample struct {
	ID        string
	SubjectID string
	// Values maps header labels to corresponding values for this row
	Values map[string]string
}

func (s Sample) GetValue(key string) (string, bool) {
	value, ok := s.Values[key]
	return value, ok
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

type Header struct {
	labels []string
	// dbGapStudyIDIndex is -1 if the header does not contain a 'dbGaP study ID' label
	dbGapStudyIDIndex int
}

func (h Header) Labels() []string {
	return h.labels
}

func ParseHeaderRow(row []string) (Header, bool, error) {
	if len(row) == 0 || !utils.EqualIgnoreCaseTrimmed(row[IDIndex], IDLabel) {
		return Header{}, false, nil
	}
	if len(row) < minLength || !utils.EqualIgnoreCaseTrimmed(row[SubjectIDIndex], SubjectIDLabel) {
		return Header{}, true, fmt.Errorf("samples header: column %d must be %q", SubjectIDIndex, SubjectIDLabel)
	}

	dbGapStudyIDIndex := slices.IndexFunc(row, func(label string) bool {
		return utils.EqualIgnoreCaseTrimmed(label, DbGapStudyIDLabel)
	})

	if dbGapStudyIDIndex < 0 {
		logger.Warn("samples file does not contain dbGaP study ID label; no samples will be included")
	}

	// defensive copy that uses the canonical ID labels even if actual header had different case
	labels := slices.Clone(row)
	labels[IDIndex] = IDLabel
	labels[SubjectIDIndex] = SubjectIDLabel
	if dbGapStudyIDIndex >= 0 {
		labels[dbGapStudyIDIndex] = DbGapStudyIDLabel
	}

	return Header{labels: labels, dbGapStudyIDIndex: dbGapStudyIDIndex}, true, nil
}

// FromRow converts the given non-header row to a Sample.
// Blank rows and rows with empty "dbGaP study ID" values are skipped: for those, a nil Sample and a nil error are
// returned so that the caller skips the row without failing.
func FromRow(header Header, row []string) (*Sample, error) {
	if utils.IsBlankRow(row) {
		logger.Info("skipping blank samples row")
		return nil, nil
	}
	if len(row) < minLength {
		return nil, fmt.Errorf("samples row is too short to contain sample and subject ids")
	}

	// Samples without a "dbGaP study ID" are intermediate derivations and should not be submitted to dbGaP.
	// Since excelize does not give us empty cells beyond the last non-empty cell,
	// dbGapStudyIDIndex >= len(row) means blank dbGap study ID.
	if header.dbGapStudyIDIndex < 0 || header.dbGapStudyIDIndex >= len(row) || utils.IsBlankCell(row[header.dbGapStudyIDIndex]) {
		logger.Info("skipping sample with no dbGaP study ID", slog.String("sampleId", row[IDIndex]))
		return nil, nil
	}

	values := make(map[string]string, len(row)-2)

	sample := Sample{
		ID:        row[IDIndex],
		SubjectID: row[SubjectIDIndex],
		Values:    values,
	}

	for i, label := range header.Labels() {
		if i >= len(row) {
			break // excelize does not give us empty cells beyond the last non-empty cell
		}
		switch i {
		case IDIndex, SubjectIDIndex:
			// skip these, they are already in struct
		default:
			values[label] = row[i]
		}
	}
	logger.Info("found sample", sample.LogGroup())
	return &sample, nil
}

func FromFile(file *excelize.File) ([]string, []Sample, error) {
	return utils.FromFile(file, ParseHeaderRow, FromRow)
}
