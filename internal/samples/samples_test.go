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

	require.Len(t, samps, 11)

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

// newSamplesFile returns an in-memory Excel file whose first sheet contains rows.
func newSamplesFile(t *testing.T, rows [][]string) *excelize.File {
	t.Helper()
	file := excelize.NewFile()
	t.Cleanup(func() {
		assert.NoError(t, file.Close())
	})
	sheet := file.GetSheetList()[0]
	for i, row := range rows {
		cell, err := excelize.CoordinatesToCellName(1, i+1)
		require.NoError(t, err)
		require.NoError(t, file.SetSheetRow(sheet, cell, &row))
	}
	return file
}

func TestFromFile_SamplesWithoutDbGapStudyIDAreSkipped(t *testing.T) {
	file := newSamplesFile(t, [][]string{
		{"Sample ID", "Subject ID", " DBGAP study id ", "nCells"},
		{"sam-1", "sub-1", "phs.111.222", "42"},
		{"sam-2", "sub-1", "", "43"},
		{"sam-3", "sub-2", " ", "44"},
		{"sam-4", "sub-2", "phs.111.222"},
		// no cells beyond the subject id, so the study id is past the end of the row
		{"sam-5", "sub-3"},
	})

	header, samps, err := FromFile(file)
	require.NoError(t, err)

	assert.Equal(t, []string{IDLabel, SubjectIDLabel, DbGapStudyIDLabel, "nCells"}, header)

	ids := make([]string, 0, len(samps))
	for _, s := range samps {
		ids = append(ids, s.ID)
	}
	assert.Equal(t, []string{"sam-1", "sam-4"}, ids)
	assert.Equal(t, "phs.111.222", samps[0].Values[DbGapStudyIDLabel])
}

func TestFromFile_NoDbGapStudyIDColumn(t *testing.T) {
	file := newSamplesFile(t, [][]string{
		{IDLabel, SubjectIDLabel, "nCells"},
		{"sam-1", "sub-1", "42"},
	})

	header, samps, err := FromFile(file)
	require.NoError(t, err)

	assert.Equal(t, []string{IDLabel, SubjectIDLabel, "nCells"}, header)
	assert.Empty(t, samps)
}

func TestParseHeaderRow(t *testing.T) {
	for name, testCase := range map[string]struct {
		row                       []string
		expectedLabels            []string
		expectedDbGapStudyIDIndex int
	}{
		"canonical labels": {
			row:                       []string{IDLabel, SubjectIDLabel, "nCells", DbGapStudyIDLabel},
			expectedLabels:            []string{IDLabel, SubjectIDLabel, "nCells", DbGapStudyIDLabel},
			expectedDbGapStudyIDIndex: 3,
		},
		"non-canonical labels": {
			row:                       []string{" Sample ID", "SUBJECT ID ", "DBGAP Study Id", "nCells"},
			expectedLabels:            []string{IDLabel, SubjectIDLabel, DbGapStudyIDLabel, "nCells"},
			expectedDbGapStudyIDIndex: 2,
		},
		"no dbGaP study ID": {
			row:                       []string{IDLabel, SubjectIDLabel, "nCells"},
			expectedLabels:            []string{IDLabel, SubjectIDLabel, "nCells"},
			expectedDbGapStudyIDIndex: -1,
		},
	} {
		t.Run(name, func(t *testing.T) {
			header, isHeader, err := ParseHeaderRow(testCase.row)
			require.NoError(t, err)
			require.True(t, isHeader)
			assert.Equal(t, testCase.expectedLabels, header.Labels())
			assert.Equal(t, testCase.expectedDbGapStudyIDIndex, header.dbGapStudyIDIndex)
		})
	}
}

func TestParseHeaderRow_LabelsAreACopy(t *testing.T) {
	row := []string{"Sample ID", "Subject ID", "nCells"}

	header, isHeader, err := ParseHeaderRow(row)
	require.NoError(t, err)
	require.True(t, isHeader)

	header.Labels()[2] = "changed"
	assert.Equal(t, []string{"Sample ID", "Subject ID", "nCells"}, row)
}

func TestParseHeaderRow_NotAHeader(t *testing.T) {
	for name, row := range map[string][]string{
		"no cells": {},
		"data row": {"sam-1", "sub-1", "phs.111.222"},
		// the sample id label is only recognized in the first column
		"sample id not first": {SubjectIDLabel, IDLabel},
	} {
		t.Run(name, func(t *testing.T) {
			_, isHeader, err := ParseHeaderRow(row)
			require.NoError(t, err)
			assert.False(t, isHeader)
		})
	}
}

func TestParseHeaderRow_MalformedHeader(t *testing.T) {
	for name, row := range map[string][]string{
		"no subject id column":    {IDLabel},
		"wrong subject id column": {IDLabel, "nCells", SubjectIDLabel},
	} {
		t.Run(name, func(t *testing.T) {
			_, isHeader, err := ParseHeaderRow(row)
			assert.Error(t, err)
			assert.True(t, isHeader)
		})
	}
}

func TestFromRow(t *testing.T) {
	header := Header{
		labels:            []string{IDLabel, SubjectIDLabel, "nCells", DbGapStudyIDLabel},
		dbGapStudyIDIndex: 3,
	}

	sample, err := FromRow(header, []string{"sam-1", "sub-1", "42", "phs.111.222"})
	require.NoError(t, err)

	require.NotNil(t, sample)
	assert.Equal(t, "sam-1", sample.ID)
	assert.Equal(t, "sub-1", sample.SubjectID)
	assert.Equal(t, map[string]string{"nCells": "42", DbGapStudyIDLabel: "phs.111.222"}, sample.Values)
}

func TestFromRow_BlankRowsAreSkipped(t *testing.T) {
	header := Header{
		labels:            []string{IDLabel, SubjectIDLabel, "nCells", DbGapStudyIDLabel},
		dbGapStudyIDIndex: 3,
	}

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

func TestFromRow_NoDbGapStudyIDIsSkipped(t *testing.T) {
	labels := []string{IDLabel, SubjectIDLabel, DbGapStudyIDLabel, "nCells"}
	withStudyID := Header{labels: labels, dbGapStudyIDIndex: 2}
	withoutStudyID := Header{labels: labels[:2], dbGapStudyIDIndex: -1}

	for name, testCase := range map[string]struct {
		header Header
		row    []string
	}{
		"empty study id":                  {header: withStudyID, row: []string{"sam-1", "sub-1", "", "42"}},
		"blank study id":                  {header: withStudyID, row: []string{"sam-1", "sub-1", " \t", "42"}},
		"study id past end":               {header: withStudyID, row: []string{"sam-1", "sub-1"}},
		"no study id column":              {header: withoutStudyID, row: []string{"sam-1", "sub-1"}},
		"no study id column, extra cells": {header: withoutStudyID, row: []string{"sam-1", "sub-1", "phs.111.222"}},
	} {
		t.Run(name, func(t *testing.T) {
			sample, err := FromRow(testCase.header, testCase.row)
			require.NoError(t, err)
			assert.Nil(t, sample)
		})
	}
}

func TestFromRow_Errors(t *testing.T) {
	header := Header{
		labels:            []string{IDLabel, SubjectIDLabel, "nCells", DbGapStudyIDLabel},
		dbGapStudyIDIndex: 3,
	}

	for name, row := range map[string][]string{
		"no subject id cell": {"sam-1"},
	} {
		t.Run(name, func(t *testing.T) {
			sample, err := FromRow(header, row)
			assert.Error(t, err)
			assert.Nil(t, sample)
		})
	}
}
