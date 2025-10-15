package subjectsample

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectsample/ds"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
	"path/filepath"
	"testing"
)

func TestWriteFiles(t *testing.T) {

	consentedSamples := []samples.Sample{
		{ID: "samp-1-sub-1", SubjectID: "sub-1"},
		{ID: "samp-1-sub-3", SubjectID: "sub-3"},
		{ID: "samp-2-sub-3", SubjectID: "sub-3"},
		{ID: "samp-3-sub-3", SubjectID: "sub-3"},
	}

	outputDir := t.TempDir()

	require.NoError(t, WriteFiles(outputDir, consentedSamples))

	actualDSFile, err := excelize.OpenFile(filepath.Join(outputDir, ds.Spec.FileName))
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, actualDSFile.Close())
	}()

	actualRows, err := actualDSFile.GetRows("Sheet1")
	require.NoError(t, err)

	// 1 Header + 4 data rows == 5 total rows
	assert.Len(t, actualRows, 5)

	assert.Equal(t, []string{dd.SubjectIDVar.Name, dd.SampleIDVar.Name}, actualRows[0])

	for i, sample := range consentedSamples {
		assert.Equal(t, []string{sample.SubjectID, sample.ID}, actualRows[i+1])
	}

}
