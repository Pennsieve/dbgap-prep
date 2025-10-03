package subjects

import (
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
	subs, err := FromFile(file)
	require.NoError(t, err)

	require.Len(t, subs, 5)

	assert.Equal(t, "sub-111", subs[0].ID)
	assert.Equal(t, "1", subs[0].Sex)

	assert.Equal(t, "sub-222", subs[1].ID)
	assert.Equal(t, "2", subs[1].Sex)

	assert.Equal(t, "sub-abc", subs[2].ID)
	assert.Equal(t, "3", subs[2].Sex)

	assert.Equal(t, "sub-xyz", subs[3].ID)
	assert.Equal(t, "4", subs[3].Sex)

	assert.Equal(t, "sub-a1b2", subs[4].ID)
	assert.Equal(t, "", subs[4].Sex)
}
