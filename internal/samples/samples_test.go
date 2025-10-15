package samples

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
	"path/filepath"
	"testing"
)

func TestFromFile(t *testing.T) {
	path := filepath.Join("testdata", FileName)
	file, err := excelize.OpenFile(path)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, file.Close())
	}()
	samplesHeader, samps, err := FromFile(file)
	require.NoError(t, err)

	fmt.Println(samplesHeader)

	assert.Equal(t, IDLabel, samplesHeader[IDIndex])
	assert.Equal(t, SubjectIDLabel, samplesHeader[SubjectIDIndex])

	assert.Len(t, samps, 11)

	sampleID := "sam-hSG-lib-8"

	assert.Equal(t, sampleID, samps[10].ID)
	expectedLabel := "nFeature_RNA"
	require.Contains(t, samps[10].Values, expectedLabel)
	// This is what Excel will show for the cell in the test file.
	// Numbers and Google Sheets will display a different value.
	// May have to re-visit based on user feedback.
	expectedValue := "1712.715569"
	assert.Equal(t, expectedValue, samps[10].Values[expectedLabel])

}
