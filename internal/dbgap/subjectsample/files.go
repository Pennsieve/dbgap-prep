package subjectsample

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	ssmdd "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectsample/dd"
	ssmds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectsample/ds"
	"github.com/pennsieve/dbgap-prep/internal/samples"
)

func WriteFiles(outputDirectory string, consentedSubjectSamples []samples.Sample) error {
	ddWriter := dd.NewNoOpWriter(outputDirectory, ssmdd.Spec.FileName)

	if err := ddWriter.Write(ssmdd.Spec); err != nil {
		return fmt.Errorf("error writing subject sample mapping file: %w", err)
	}

	dsWriter := ds.NewXLSXWriter(outputDirectory, ssmds.DefaultFileNameBase)
	if err := ssmds.Write(dsWriter, consentedSubjectSamples); err != nil {
		return fmt.Errorf("error writing subject sample mapping file: %w", err)
	}

	return nil
}
