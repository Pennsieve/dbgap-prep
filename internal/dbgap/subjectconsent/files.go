package subjectconsent

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	scdd "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/dd"
	scds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/ds"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
	"log/slog"
	"path/filepath"
)

var logger = logging.PackageLogger("subjectconsent")

func WriteFiles(outputDirectory string, subs []subjects.Subject) ([]scds.SubjectConsent, error) {
	subjectConsentDDPath := filepath.Join(outputDirectory, scdd.Spec.FileName)

	if err := dd.Write(subjectConsentDDPath, scdd.Spec); err != nil {
		return nil, err
	}

	logger.Info("wrote subject consent DD file", slog.String("file", subjectConsentDDPath))

	subjectConsentDSPath := filepath.Join(outputDirectory, scds.Spec.FileName)
	subjectConsents, err := scds.Write(subjectConsentDSPath, subs)
	if err != nil {
		return nil, err
	}

	logger.Info("wrote subject consent DS file",
		slog.String("file", subjectConsentDSPath),
		slog.Int("subjectConsentCount", len(subjectConsents)))

	return subjectConsents, nil
}
