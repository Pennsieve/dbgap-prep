package subjectsample

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	ssmdd "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectsample/dd"
	ssmds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectsample/ds"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"log/slog"
	"path/filepath"
)

var logger = logging.PackageLogger("subjectsample")

func WriteFiles(outputDirectory string, consentedSubjectSamples map[string][]samples.Sample) error {
	subjectSampleMappingDDPath := filepath.Join(outputDirectory, ssmdd.Spec.FileName)

	if err := dd.Write(subjectSampleMappingDDPath, ssmdd.Spec); err != nil {
		return err
	}

	logger.Info("wrote subject sample mapping DD file", slog.String("file", subjectSampleMappingDDPath))

	subjectToSamples := make(map[string][]string, len(consentedSubjectSamples))
	for subjectID, samps := range consentedSubjectSamples {
		sampleIDs := make([]string, len(samps))
		for i, sample := range samps {
			sampleIDs[i] = sample.ID
		}
		subjectToSamples[subjectID] = sampleIDs
	}

	subjectSampleMappingDSPath := filepath.Join(outputDirectory, ssmds.Spec.FileName)
	if err := ssmds.Write(subjectSampleMappingDSPath, subjectToSamples); err != nil {
		return err
	}

	logger.Info("wrote subject sample mapping DS file",
		slog.String("file", subjectSampleMappingDSPath),
		slog.Int("subjectSampleCount", len(subjectToSamples)))
	return nil
}
