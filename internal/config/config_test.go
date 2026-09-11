package config

import (
	"testing"

	"github.com/pennsieve/dbgap-prep/internal/enums/analytetype"
	"github.com/pennsieve/dbgap-prep/internal/enums/consentgroup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyteType_String(t *testing.T) {
	assert.Equal(t, "DNA", analytetype.DNA.String())
	assert.Equal(t, "RNA", analytetype.RNA.String())
	assert.Equal(t, "DNA/RNA", analytetype.DNARNA.String())
	assert.Equal(t, "UNKNOWN", analytetype.Type(-1).String())
}

func TestAnalyteTypeFromString(t *testing.T) {
	testCases := map[string]analytetype.Type{
		"DNA":     analytetype.DNA,
		"dna":     analytetype.DNA,
		"RNA":     analytetype.RNA,
		"rna":     analytetype.RNA,
		"DNA/RNA": analytetype.DNARNA,
		"dna/rna": analytetype.DNARNA,
	}
	for input, expected := range testCases {
		t.Run(input, func(t *testing.T) {
			actual, err := analytetype.FromString(input)
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestAnalyteTypeFromString_Invalid(t *testing.T) {
	_, err := analytetype.FromString("not-a-real-analyte-type")
	assert.Error(t, err)
}

func TestConsentGroup_String(t *testing.T) {
	assert.Equal(t, "GRU", consentgroup.GRU.String())
	assert.Equal(t, "HMB", consentgroup.HMB.String())
	assert.Equal(t, "UNKNOWN", consentgroup.Group(-1).String())
}

func TestConsentGroupFromString(t *testing.T) {
	testCases := map[string]consentgroup.Group{
		"GRU": consentgroup.GRU,
		"gru": consentgroup.GRU,
		"HMB": consentgroup.HMB,
		"hmb": consentgroup.HMB,
	}
	for input, expected := range testCases {
		t.Run(input, func(t *testing.T) {
			actual, err := consentgroup.FromString(input)
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestConsentGroupFromString_Invalid(t *testing.T) {
	_, err := consentgroup.FromString("not-a-real-consent-group")
	assert.Error(t, err)
}
