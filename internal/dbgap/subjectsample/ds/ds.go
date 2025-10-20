package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	"github.com/pennsieve/dbgap-prep/internal/samples"
)

const DefaultFileNameBase = "3a_SSM_DS"

var Spec = ds.Spec{
	Variables: []dd.Variable{*dd.SubjectIDVar, *dd.SampleIDVar},
}

func Write(writer ds.Writer, subjectSamples []samples.Sample) error {
	var rows [][]string

	for _, sample := range subjectSamples {
		row := []string{sample.SubjectID, sample.ID}
		rows = append(rows, row)
	}

	return writer.Write(Spec, rows)
}
