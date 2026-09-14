package subjects

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
	header, subs, err := FromFile(file)
	require.NoError(t, err)

	assert.Equal(t, IDLabel, header[IDIndex])

	require.Len(t, subs, 5)

	assert.Equal(t, "sub-111", subs[0].ID)
	assert.Equal(t, "1", subs[0].Sex)
	assert.Equal(t, "dbgap#111", subs[0].SourceSubjectIDdbGap)

	assert.Equal(t, "sub-222", subs[1].ID)
	assert.Equal(t, "2", subs[1].Sex)
	assert.Equal(t, "dbgap#222", subs[1].SourceSubjectIDdbGap)

	assert.Equal(t, "sub-abc", subs[2].ID)
	assert.Equal(t, "3", subs[2].Sex)
	assert.Equal(t, "dbgap#abc", subs[2].SourceSubjectIDdbGap)

	assert.Equal(t, "sub-xyz", subs[3].ID)
	assert.Equal(t, "4", subs[3].Sex)
	assert.Equal(t, "dbgap#xyz", subs[3].SourceSubjectIDdbGap)

	assert.Equal(t, "sub-a1b2", subs[4].ID)
	assert.Equal(t, "", subs[4].Sex)
	assert.Equal(t, "dbgap#a1b2", subs[4].SourceSubjectIDdbGap)

}

func TestFromRow(t *testing.T) {
	header := []string{IDLabel, SexLabel, SourceSubjectIDdbGapColumn, "age"}

	// An empty expectedID means the row is expected to be skipped.
	for name, testCase := range map[string]struct {
		row        []string
		expectedID string
	}{
		"subject row":            {row: []string{"sub-1", "1", "dbgap#1", "42"}, expectedID: "sub-1"},
		"uppercase prefix":       {row: []string{"SUB-1", "1", "dbgap#1", "42"}, expectedID: "SUB-1"},
		"untrimmed prefix":       {row: []string{" sub-1", "1", "dbgap#1", "42"}, expectedID: " sub-1"},
		"non-subject row":        {row: []string{"1234", "1", "dbgap#1", "42"}},
		"empty id":               {row: []string{"", "", "", ""}},
		"blank id":               {row: []string{"   ", "", "", ""}},
		"id shorter than prefix": {row: []string{"sub", "", "", ""}},
		"blank row":              {row: []string{}},
		"blank row with cells":   {row: []string{"", " "}},
	} {
		t.Run(name, func(t *testing.T) {
			subject, err := FromRow(header, testCase.row)
			require.NoError(t, err)

			if len(testCase.expectedID) == 0 {
				assert.Nil(t, subject)
				return
			}

			require.NotNil(t, subject)
			assert.Equal(t, testCase.expectedID, subject.ID)
			assert.Equal(t, "1", subject.Sex)
			assert.Equal(t, "dbgap#1", subject.SourceSubjectIDdbGap)
			assert.Equal(t, map[string]string{"age": "42"}, subject.Values)
		})
	}
}

func TestFromRow_HeaderRowIsError(t *testing.T) {
	header := []string{IDLabel, SexLabel}

	subject, err := FromRow(header, []string{IDLabel, SexLabel})

	assert.Error(t, err)
	assert.Nil(t, subject)
}
