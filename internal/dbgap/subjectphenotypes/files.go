package subjectphenotypes

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	subjectgphenotypesdd "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectphenotypes/dd"
	subjectgphenotypesds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectphenotypes/ds"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
	"log/slog"
	"path/filepath"
)

var logger = logging.PackageLogger("subjectphenotypes")

func WriteFiles(outputDirectory string, subjectsHeader []string, consentedSubjects []subjects.Subject) error {
	attrLabels := HeaderToAttributeLabels(subjectsHeader)
	variables := subjectgphenotypesdd.Variables(attrLabels)
	spec := subjectgphenotypesdd.Spec(variables)
	subjectPhenotypesDDPath := filepath.Join(outputDirectory, spec.FileName)

	if err := dd.Write(subjectPhenotypesDDPath, spec); err != nil {
		return err
	}

	logger.Info("wrote subject phenotypes DD file", slog.String("file", subjectPhenotypesDDPath))

	subjectPhenotypesDSPath := filepath.Join(outputDirectory, subjectgphenotypesds.FileName)
	if err := subjectgphenotypesds.Write(subjectPhenotypesDSPath, variables, consentedSubjects); err != nil {
		return err
	}

	logger.Info("wrote subject phenotypes DS file", slog.String("file", subjectPhenotypesDSPath))

	return nil
}

func HeaderToAttributeLabels(subjectsHeader []string) []string {
	cleaned := make([]string, 0, len(subjectsHeader))
	for _, label := range subjectsHeader {
		if label != subjects.IDLabel && label != subjects.SexLabel {
			cleaned = append(cleaned, label)
		}
	}
	return cleaned
}
