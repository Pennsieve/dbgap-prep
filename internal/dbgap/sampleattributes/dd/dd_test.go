package dd

import (
	"testing"

	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/models"
	"github.com/stretchr/testify/assert"
)

func TestSpec_EncodedValueVariable(t *testing.T) {
	spec := Spec([]dd.Variable{models.IsTumorVar})

	require := assert.New(t)
	require.Equal("6b_SampleAttributes_DD.xlsx", spec.FileName)
	require.Len(spec.Rows, 1)

	row := spec.Rows[0]
	// header is VARNAME, VARDESC, TYPE, VALUES; VALUES expands to one cell per encoded value
	require.Equal(models.IsTumorVar.Name, row[0])
	require.Equal(models.IsTumorVar.Description, row[1])
	require.Equal(models.IsTumorVar.Type, row[2])
	require.Equal("Y=Is tumor", row[3].(dd.EncodedValue).String())
	require.Equal("N=Is not a tumor", row[4].(dd.EncodedValue).String())
	require.Equal("UNK=Tumor status unknown", row[5].(dd.EncodedValue).String())
}

func TestSpec_StringVariable(t *testing.T) {
	spec := Spec([]dd.Variable{models.AnalyteTypeVar})

	row := spec.Rows[0]
	// no VALUES since AnalyteTypeVar has none
	assert.Len(t, row, 3)
	assert.Equal(t, models.AnalyteTypeVar.Name, row[0])
	assert.Equal(t, models.AnalyteTypeVar.Description, row[1])
	assert.Equal(t, models.AnalyteTypeVar.Type, row[2])
}
