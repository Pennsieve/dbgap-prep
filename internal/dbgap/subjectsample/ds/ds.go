package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
)

var Spec = ds.Spec{
	FileName: "3a_SSM_DS.txt",
	Header:   []string{dd.SubjectIDVar.Name, dd.SampleIDVar.Name},
}

func Write(path string, subjectToSamples map[string][]string) error {
	var rows [][]string

	for subjectID, sampleIDs := range subjectToSamples {
		rowsForSubject := make([][]string, 0, len(sampleIDs))
		for _, sampleID := range sampleIDs {
			row := []string{subjectID, sampleID}
			rowsForSubject = append(rowsForSubject, row)
		}
		rows = append(rows, rowsForSubject...)
	}

	return ds.Write(path, Spec, rows)
}
