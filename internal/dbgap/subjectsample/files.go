package subjectsample

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	scds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/ds"
	ssmdd "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectsample/dd"
	ssmds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectsample/ds"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"log/slog"
	"path/filepath"
)

var logger = logging.PackageLogger("subjectsample")

func WriteFiles(outputDirectory string, subjectConsents []scds.SubjectConsent, samps []samples.Sample) error {
	subjectSampleMappingDDPath := filepath.Join(outputDirectory, ssmdd.Spec.FileName)

	if err := dd.Write(subjectSampleMappingDDPath, ssmdd.Spec); err != nil {
		return err
	}

	logger.Info("wrote subject sample mapping DD file", slog.String("file", subjectSampleMappingDDPath))

	subjectToSamples := make(map[string][]string, len(subjectConsents))
	for _, sample := range samps {
		// Don't include samples that have no subject?
		if sample.HasSubject() {
			subjectToSamples[sample.SubjectID] = append(subjectToSamples[sample.SubjectID], sample.ID)
		}
	}

	// Remove the subjects with no consent
	for _, subjectConsent := range subjectConsents {
		if !subjectConsent.IsConsented() {
			delete(subjectToSamples, subjectConsent.SubjectID)
		}
	}

	subjectSampleMappingDSPath := filepath.Join(outputDirectory, ssmds.Spec.FileName)
	if err := ssmds.Write(subjectSampleMappingDSPath, subjectToSamples); err != nil {
		return err
	}

	logger.Info("wrote subject sample mapping DS file", slog.String("file", subjectSampleMappingDSPath))
	return nil
}
