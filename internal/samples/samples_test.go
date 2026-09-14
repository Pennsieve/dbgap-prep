package samples

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
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

func TestFromRow(t *testing.T) {
	header := []string{IDLabel, SubjectIDLabel, "nCells"}

	sample, err := FromRow(header, []string{"sam-1", "sub-1", "42"})
	require.NoError(t, err)

	require.NotNil(t, sample)
	assert.Equal(t, "sam-1", sample.ID)
	assert.Equal(t, "sub-1", sample.SubjectID)
	assert.Equal(t, map[string]string{"nCells": "42"}, sample.Values)
}

func TestFromRow_BlankRowsAreSkipped(t *testing.T) {
	header := []string{IDLabel, SubjectIDLabel, "nCells"}

	for name, row := range map[string][]string{
		"no cells":    {},
		"empty cells": {"", "", ""},
		"blank cells": {" ", "", "\t"},
	} {
		t.Run(name, func(t *testing.T) {
			sample, err := FromRow(header, row)
			require.NoError(t, err)
			assert.Nil(t, sample)
		})
	}
}

func TestFromRow_Errors(t *testing.T) {
	header := []string{IDLabel, SubjectIDLabel, "nCells"}

	for name, row := range map[string][]string{
		"header row":         {IDLabel, SubjectIDLabel, "nCells"},
		"no subject id cell": {"sam-1"},
	} {
		t.Run(name, func(t *testing.T) {
			sample, err := FromRow(header, row)
			assert.Error(t, err)
			assert.Nil(t, sample)
		})
	}
}
