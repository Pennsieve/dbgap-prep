package subjects

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/logging"
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
}

func (s Subject) String() string {
	return fmt.Sprintf("subject: id = [%s], sex = [%s]", s.ID, s.Sex)
}

func (s Subject) LogGroup() slog.Attr {
	return slog.Group("subject", slog.String("id", s.ID), slog.String("sex", s.Sex))
}

func IsHeaderRow(row []string) bool {
	return len(row) > 0 &&
		row[IDIndex] == IDLabel
}

func FromRow(row []string) (*Subject, error) {
	// skip empty rows and header rows.
	// The excelize library shouldn't really give us empty rows,
	// but the id in first column is required, so make sure we
	// don't panic below.
	if len(row) == 0 || IsHeaderRow(row) {
		return nil, nil
	}

	subject := Subject{ID: row[IDIndex]}
	if len(row) > SexIndex {
		subject.Sex = row[SexIndex]
	}
	logger.Info("found subject", subject.LogGroup())
	return &subject, nil

}

// FromRows returns a map from subject id to Subject of all subjects found in rows.
func FromRows(rows [][]string) ([]Subject, error) {
	allSubs := make([]Subject, 0, len(rows))
	for _, row := range rows {
		if s, err := FromRow(row); err != nil {
			return nil, err
		} else if s != nil {
			allSubs = append(allSubs, *s)
		}
	}
	return allSubs, nil
}

func FromFile(subjectsFile *excelize.File) ([]Subject, error) {
	var allSubs []Subject
	for _, subjectSheet := range subjectsFile.GetSheetList() {
		rows, err := subjectsFile.GetRows(subjectSheet)
		if err != nil {
			return nil, fmt.Errorf("error getting rows from %s, sheet %s: %w", subjectsFile.Path, subjectSheet, err)
		}

		subs, err := FromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("error getting subjects from %s, sheet %s: %w", subjectsFile.Path, subjectSheet, err)
		}
		allSubs = append(allSubs, subs...)
	}
	return allSubs, nil
}
