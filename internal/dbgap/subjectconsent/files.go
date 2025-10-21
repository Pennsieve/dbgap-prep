package subjectconsent

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	scdd "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/dd"
	scds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/ds"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
	"log/slog"
)

var logger = logging.PackageLogger("subjectconsent")

func WriteFiles(outputDirectory string, subs []subjects.Subject) ([]scds.SubjectConsent, error) {
	ddWriter := dd.NewNoOpWriter(outputDirectory, scdd.Spec.FileName)

	if err := ddWriter.Write(scdd.Spec); err != nil {
		return nil, err
	}

	dsWriter := ds.NewXLSXWriter(outputDirectory, scds.DefaultFileNameBase)
	subjectConsents, err := scds.Write(dsWriter, subs)
	if err != nil {
		return nil, err
	}

	logger.Info("got subject consents",
		slog.Int("count", len(subjectConsents)))

	return subjectConsents, nil
}
