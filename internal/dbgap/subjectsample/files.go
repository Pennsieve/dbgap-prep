package subjectsample

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	ssmdd "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectsample/dd"
	ssmds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectsample/ds"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"log/slog"
	"path/filepath"
)

var logger = logging.PackageLogger("subjectsample")

func WriteFiles(outputDirectory string, consentedSubjectSamples []samples.Sample) error {
	subjectSampleMappingDDPath := filepath.Join(outputDirectory, ssmdd.Spec.FileName)

	if err := dd.Write(subjectSampleMappingDDPath, ssmdd.Spec); err != nil {
		return err
	}

	logger.Info("wrote subject sample mapping DD file", slog.String("file", subjectSampleMappingDDPath))

	dsWriter := ds.NewXLSXWriter(outputDirectory, ssmds.DefaultFileNameBase)
	if err := ssmds.Write(dsWriter, consentedSubjectSamples); err != nil {
		return err
	}

	logger.Info("wrote subject sample mapping DS file",
		slog.String("file", dsWriter.Path()),
	)
	return nil
}
