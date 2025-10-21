package ds

import (
	"github.com/google/uuid"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
	"strings"
	"testing"
)

func TestXLSXWriter_WriteWideColumns(t *testing.T) {
	writer := NewXLSXWriter(t.TempDir(), "test-ds")

	expectedIDVariable := dd.Variable{Name: "id", Description: "the id", Type: dd.StringType}
	expectedValueVariable := dd.Variable{Name: "value", Description: "the value", Type: dd.StringType}
	spec := Spec{Variables: []dd.Variable{expectedIDVariable, expectedValueVariable}}

	rows := [][]string{
		{uuid.NewString(), strings.Repeat("1", 300)},
		{uuid.NewString(), uuid.NewString()},
	}

	require.NoError(t, writer.Write(spec, rows))

	actualDSFile, err := excelize.OpenFile(writer.Path())
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, actualDSFile.Close())
	}()

	actualRows, err := actualDSFile.GetRows("Sheet1")
	require.NoError(t, err)

	require.Len(t, actualRows, len(rows)+1)

	//header
	header := actualRows[0]
	assert.Equal(t, []string{expectedIDVariable.Name, expectedValueVariable.Name}, header)

	//data
	assert.Equal(t, rows, actualRows[1:])
}
