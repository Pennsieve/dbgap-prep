package subjectconsent

import (
	"fmt"
	"log/slog"

	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	scdd "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/dd"
	scds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/ds"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/models"
	"github.com/pennsieve/dbgap-prep/internal/enums/consentgroup"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

var logger = logging.PackageLogger("subjectconsent")

func WriteFiles(outputDirectory string, consentGroup consentgroup.Group, subs []subjects.Subject) ([]scds.SubjectConsent, error) {
	consentVariable, err := models.ConsentVariable(consentGroup)
	if err != nil {
		return nil, err
	}
	ddSpec := scdd.Spec(&consentVariable)
	ddWriter := dd.NewXLSXWriter(outputDirectory, ddSpec.FileName)

	if err := ddWriter.Write(ddSpec); err != nil {
		return nil, fmt.Errorf("error writing subject consent file: %w", err)
	}

	dsWriter := ds.NewXLSXWriter(outputDirectory, scds.DefaultFileNameBase)
	subjectConsents, err := scds.Write(dsWriter, scds.Spec(consentVariable), subs)
	if err != nil {
		return nil, fmt.Errorf("error writing subject consent file: %w", err)
	}

	logger.Info("got subject consents",
		slog.Int("count", len(subjectConsents)))

	return subjectConsents, nil
}
