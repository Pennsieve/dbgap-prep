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

func TestLookupIsTumorEnvVar(t *testing.T) {
	const key = "IS_TUMOR"

	t.Run("unset uses default", func(t *testing.T) {
		t.Setenv(key, "")
		value, err := LookupIsTumorEnvVar(key, config.YES)
		require.NoError(t, err)
		assert.Equal(t, config.YES, value)

		value, err = LookupIsTumorEnvVar(key, config.NO)
		require.NoError(t, err)
		assert.Equal(t, config.NO, value)
	})

	for _, yes := range []string{"yes", "YES", "Yes"} {
		t.Run("yes/"+yes, func(t *testing.T) {
			t.Setenv(key, yes)
			value, err := LookupIsTumorEnvVar(key, config.NO)
			require.NoError(t, err)
			assert.Equal(t, config.YES, value)
		})
	}

	for _, no := range []string{"no", "NO", "No"} {
		t.Run("no/"+no, func(t *testing.T) {
			t.Setenv(key, no)
			value, err := LookupIsTumorEnvVar(key, config.YES)
			require.NoError(t, err)
			assert.Equal(t, config.NO, value)
		})
	}

	t.Run("unparsable errors instead of using default", func(t *testing.T) {
		t.Setenv(key, "not-an-is-tumor")
		_, err := LookupIsTumorEnvVar(key, config.NO)
		assert.Error(t, err)
	})
}
