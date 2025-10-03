package dd

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVariable_ToDDRow_EncodedValues(t *testing.T) {
	variable := Variable{
		Name:        uuid.NewString(),
		Description: uuid.NewString(),
		Type:        EncodedValueType,
		Values: []EncodedValue{
			NewEncodedValue("1", uuid.NewString()),
			NewEncodedValue("2", uuid.NewString()),
			NewEncodedValue("OTHER", uuid.NewString()),
		},
	}

	ddRow := variable.ToDDRow()

	// should have elements for name, description, type, plus one for each encoded
	require.Len(t, ddRow, 3+len(variable.Values))
	assert.Equal(t, variable.Name, ddRow[0])
	assert.Equal(t, variable.Description, ddRow[1])
	assert.Equal(t, variable.Type, ddRow[2])
	for i, expectedValue := range variable.Values {
		actual := ddRow[3+i]
		var actualValue EncodedValue
		require.IsType(t, actualValue, actual)
		actualValue = actual.(EncodedValue)

		assert.Equal(t, expectedValue.String(), actualValue.String())
	}
}

func TestVariable_ToDDRow_NotEncodedValues(t *testing.T) {
	variable := Variable{
		Name:        uuid.NewString(),
		Description: uuid.NewString(),
		Type:        StringType,
	}

	ddRow := variable.ToDDRow()

	// should have elements for name, description, type only
	require.Len(t, ddRow, 3)
	assert.Equal(t, variable.Name, ddRow[0])
	assert.Equal(t, variable.Description, ddRow[1])
	assert.Equal(t, variable.Type, ddRow[2])
}
