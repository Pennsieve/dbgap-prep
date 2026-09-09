package app

import (
	"testing"

	"github.com/pennsieve/dbgap-prep/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookupAnalyteTypeEnvVar(t *testing.T) {
	testCases := map[string]config.AnalyteType{
		"DNA":     config.DNA,
		"rna":     config.RNA,
		"DNA/RNA": config.DNARNA,
	}
	for value, expected := range testCases {
		t.Run(value, func(t *testing.T) {
			t.Setenv(AnalyteTypeKey, value)
			actual, err := LookupAnalyteTypeEnvVar()
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestLookupAnalyteTypeEnvVar_MissingOrInvalid(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		t.Setenv(AnalyteTypeKey, "")
		_, err := LookupAnalyteTypeEnvVar()
		assert.Error(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		t.Setenv(AnalyteTypeKey, "not-a-real-analyte-type")
		_, err := LookupAnalyteTypeEnvVar()
		assert.Error(t, err)
	})
}

func TestLookupBoolEnvVar(t *testing.T) {
	const key = "IS_TUMOR"

	t.Run("unset uses default", func(t *testing.T) {
		t.Setenv(key, "")
		value, err := LookupBoolEnvVar(key, true)
		require.NoError(t, err)
		assert.True(t, value)

		value, err = LookupBoolEnvVar(key, false)
		require.NoError(t, err)
		assert.False(t, value)
	})

	for _, truthy := range []string{"true", "TRUE", "1", "t"} {
		t.Run("truthy/"+truthy, func(t *testing.T) {
			t.Setenv(key, truthy)
			value, err := LookupBoolEnvVar(key, false)
			require.NoError(t, err)
			assert.True(t, value)
		})
	}

	for _, falsy := range []string{"false", "FALSE", "0", "f"} {
		t.Run("falsy/"+falsy, func(t *testing.T) {
			t.Setenv(key, falsy)
			value, err := LookupBoolEnvVar(key, true)
			require.NoError(t, err)
			assert.False(t, value)
		})
	}

	t.Run("unparsable errors instead of using default", func(t *testing.T) {
		t.Setenv(key, "not-a-bool")
		_, err := LookupBoolEnvVar(key, false)
		assert.Error(t, err)
	})
}