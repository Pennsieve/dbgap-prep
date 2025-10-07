package subjectsample

import (
	"encoding/csv"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectsample/ds"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFiles(t *testing.T) {
	consentedSamples := map[string][]samples.Sample{
		"sub-1": {
			{ID: "samp-1-sub-1", SubjectID: "sub-1"},
		},
		"sub-3": {
			{ID: "samp-1-sub-3", SubjectID: "sub-3"},
			{ID: "samp-2-sub-3", SubjectID: "sub-3"},
			{ID: "samp-3-sub-3", SubjectID: "sub-3"}},
	}

	outputDir := t.TempDir()

	require.NoError(t, WriteFiles(outputDir, consentedSamples))

	actualDSFile, err := os.Open(filepath.Join(outputDir, ds.Spec.FileName))
	require.NoError(t, err)

	tsvReader := csv.NewReader(actualDSFile)
	tsvReader.Comma = '\t'
	actualRows, err := tsvReader.ReadAll()
	require.NoError(t, err)

	// 1 Header + 1 sub-1 row + 3 sub-3 rows == 5 total rows
	assert.Len(t, actualRows, 5)

}
