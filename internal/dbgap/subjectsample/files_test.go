package subjectsample

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	ssmds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectsample/ds"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
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

	dsWriter := ds.NewXLSXWriter(outputDir, ssmds.DefaultFileNameBase)
	require.NoError(t, ssmds.Write(dsWriter, consentedSamples))

	actualDSFile, err := excelize.OpenFile(dsWriter.Path())
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
