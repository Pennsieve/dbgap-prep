package app

import (
	"os"
	"testing"

	"github.com/pennsieve/dbgap-prep/internal/enums/analytetype"
	"github.com/pennsieve/dbgap-prep/internal/enums/consentgroup"
	"github.com/pennsieve/dbgap-prep/internal/enums/istumor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// optionalEnvVarKeys are the env vars that ConfigFromEnv does not require.
var optionalEnvVarKeys = []string{WorkflowInstanceIDKey, IsTumorKey, PHSAccessionKey}

// setEnvVars sets every required env var to a valid value and unsets the optional ones, so
// that tests do not depend on the environment they are run in. The values are chosen to
// differ from the zero value of the Config field they end up in, so that an assertion on
// one of them fails if the field is never assigned.
func setEnvVars(t *testing.T) {
	t.Helper()
	t.Setenv(IntegrationIDKey, "integration-1")
	t.Setenv(InputDirectoryKey, "/tmp/input")
	t.Setenv(OutputDirectoryKey, "/tmp/output")
	t.Setenv(ConsentGroupKey, consentgroup.HMB.String())
	t.Setenv(AnalyteTypeKey, analytetype.RNA.String())

	for _, key := range optionalEnvVarKeys {
		unsetEnvVar(t, key)
	}
}

// unsetEnvVar removes the given env var for the duration of the test. The t.Setenv call
// registers the cleanup that restores whatever value the var had before the test ran.
func unsetEnvVar(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "")
	require.NoError(t, os.Unsetenv(key))
}

func TestConfigFromEnv(t *testing.T) {
	setEnvVars(t)
	t.Setenv(WorkflowInstanceIDKey, "workflow-instance-1")
	t.Setenv(IsTumorKey, istumor.YES.String())
	t.Setenv(PHSAccessionKey, "phs003456.v1.p1")

	cfg, err := ConfigFromEnv()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "integration-1", cfg.IntegrationID)
	assert.Equal(t, "workflow-instance-1", cfg.WorkflowInstanceID)
	assert.Equal(t, "/tmp/input", cfg.InputDirectory)
	assert.Equal(t, "/tmp/output", cfg.OutputDirectory)
	assert.Equal(t, consentgroup.HMB, cfg.ConsentGroup)
	assert.Equal(t, analytetype.RNA, cfg.AnalyteType)
	assert.Equal(t, istumor.YES, cfg.IsTumor)
	assert.Equal(t, "phs003456.v1.p1", cfg.PHSAccession)
}

func TestConfigFromEnv_OptionalVarsUnset(t *testing.T) {
	setEnvVars(t)

	// IS_TUMOR falls back to a default rather than to the zero value of the field, so set
	// it to the non-default value first: if unsetting it did not change the result, the
	// default is not being applied.
	t.Setenv(IsTumorKey, istumor.YES.String())
	set, err := ConfigFromEnv()
	require.NoError(t, err)
	require.NotNil(t, set)
	require.Equal(t, istumor.YES, set.IsTumor)

	unsetEnvVar(t, IsTumorKey)

	cfg, err := ConfigFromEnv()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Empty(t, cfg.WorkflowInstanceID)
	assert.Empty(t, cfg.PHSAccession)
	assert.Equal(t, istumor.NO, cfg.IsTumor)
}

// TestConfigFromEnv_Errors covers what ConfigFromEnv itself is responsible for: reading
// each key and passing the failure on. The parsing of individual values is covered by the
// LookupXEnvVar tests below, so only one failing value per lookup is included here.
func TestConfigFromEnv_Errors(t *testing.T) {
	for name, envVar := range map[string]struct {
		key   string
		value string
	}{
		"missing integration id":   {key: IntegrationIDKey},
		"missing input directory":  {key: InputDirectoryKey},
		"missing output directory": {key: OutputDirectoryKey},
		"invalid consent group":    {key: ConsentGroupKey, value: "not-a-consent-group"},
		"invalid analyte type":     {key: AnalyteTypeKey, value: "not-an-analyte-type"},
		// IS_TUMOR is optional, but an unparsable value must not fall back to the default
		"invalid is tumor": {key: IsTumorKey, value: "not-an-is-tumor"},
	} {
		t.Run(name, func(t *testing.T) {
			setEnvVars(t)
			if len(envVar.value) == 0 {
				unsetEnvVar(t, envVar.key)
			} else {
				t.Setenv(envVar.key, envVar.value)
			}

			cfg, err := ConfigFromEnv()

			assert.Error(t, err)
			assert.Nil(t, cfg)
		})
	}
}

func TestLookupAnalyteTypeEnvVar(t *testing.T) {
	testCases := map[string]analytetype.Type{
		"DNA":     analytetype.DNA,
		"rna":     analytetype.RNA,
		"DNA/RNA": analytetype.DNARNA,
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
	t.Run("unset uses default", func(t *testing.T) {
		t.Setenv(IsTumorKey, "")
		value, err := LookupIsTumorEnvVar(istumor.YES)
		require.NoError(t, err)
		assert.Equal(t, istumor.YES, value)

		value, err = LookupIsTumorEnvVar(istumor.NO)
		require.NoError(t, err)
		assert.Equal(t, istumor.NO, value)
	})

	for _, yes := range []string{"yes", "YES", "Yes"} {
		t.Run("yes/"+yes, func(t *testing.T) {
			t.Setenv(IsTumorKey, yes)
			value, err := LookupIsTumorEnvVar(istumor.NO)
			require.NoError(t, err)
			assert.Equal(t, istumor.YES, value)
		})
	}

	for _, no := range []string{"no", "NO", "No"} {
		t.Run("no/"+no, func(t *testing.T) {
			t.Setenv(IsTumorKey, no)
			value, err := LookupIsTumorEnvVar(istumor.YES)
			require.NoError(t, err)
			assert.Equal(t, istumor.NO, value)
		})
	}

	t.Run("blank uses default", func(t *testing.T) {
		t.Setenv(IsTumorKey, "   ")
		value, err := LookupIsTumorEnvVar(istumor.YES)
		require.NoError(t, err)
		assert.Equal(t, istumor.YES, value)
	})

	t.Run("surrounding whitespace is ignored", func(t *testing.T) {
		t.Setenv(IsTumorKey, "  "+istumor.YES.String()+"  ")
		value, err := LookupIsTumorEnvVar(istumor.NO)
		require.NoError(t, err)
		assert.Equal(t, istumor.YES, value)
	})

	t.Run("unparsable errors instead of using default", func(t *testing.T) {
		t.Setenv(IsTumorKey, "not-an-is-tumor")
		_, err := LookupIsTumorEnvVar(istumor.NO)
		assert.Error(t, err)
	})
}
