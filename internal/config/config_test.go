package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyteType_String(t *testing.T) {
	assert.Equal(t, "DNA", DNA.String())
	assert.Equal(t, "RNA", RNA.String())
	assert.Equal(t, "DNA/RNA", DNARNA.String())
	assert.Equal(t, "UNKNOWN", AnalyteType(-1).String())
}

func TestAnalyteTypeFromString(t *testing.T) {
	testCases := map[string]AnalyteType{
		"DNA":     DNA,
		"dna":     DNA,
		"RNA":     RNA,
		"rna":     RNA,
		"DNA/RNA": DNARNA,
		"dna/rna": DNARNA,
	}
	for input, expected := range testCases {
		t.Run(input, func(t *testing.T) {
			actual, err := AnalyteTypeFromString(input)
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestAnalyteTypeFromString_Invalid(t *testing.T) {
	_, err := AnalyteTypeFromString("not-a-real-analyte-type")
	assert.Error(t, err)
}

func TestConsentGroup_String(t *testing.T) {
	assert.Equal(t, "GRU", GRU.String())
	assert.Equal(t, "HMB", HMB.String())
	assert.Equal(t, "UNKNOWN", ConsentGroup(-1).String())
}

func TestConsentGroupFromString(t *testing.T) {
	testCases := map[string]ConsentGroup{
		"GRU": GRU,
		"gru": GRU,
		"HMB": HMB,
		"hmb": HMB,
	}
	for input, expected := range testCases {
		t.Run(input, func(t *testing.T) {
			actual, err := ConsentGroupFromString(input)
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestConsentGroupFromString_Invalid(t *testing.T) {
	_, err := ConsentGroupFromString("not-a-real-consent-group")
	assert.Error(t, err)
}