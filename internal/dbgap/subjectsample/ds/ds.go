package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	"github.com/pennsieve/dbgap-prep/internal/samples"
)

var Spec = ds.Spec{
	FileName: "3a_SSM_DS.txt",
	Header:   []string{dd.SubjectIDVar.Name, dd.SampleIDVar.Name},
}

func Write(path string, subjectSamples []samples.Sample) error {
	var rows [][]string

	for _, sample := range subjectSamples {
		row := []string{sample.SubjectID, sample.ID}
		rows = append(rows, row)
	}

	return ds.Write(path, Spec, rows)
}
