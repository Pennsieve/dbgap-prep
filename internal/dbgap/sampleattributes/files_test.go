package sampleattributes

import (
	"testing"

	"github.com/pennsieve/dbgap-prep/internal/config"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/models"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func namesOf(variables []dd.Variable) []string {
	return dd.VariableNames(variables)
}

func TestPresentVariables_AlwaysIncludesComputedVars(t *testing.T) {
	forcedVarNames := []string{models.AnalyteTypeVar.Name, models.IsTumorVar.Name, models.SPARCDatasetDOIVar.Name}

	testCases := map[string][]string{
		"empty header":            {},
		"unrelated header only":   {"some other column", "another column"},
		"missing optional source": {"sample id"},
	}

	for name, header := range testCases {
		t.Run(name, func(t *testing.T) {
			present := namesOf(presentVariables(header))
			for _, forced := range forcedVarNames {
				assert.Contains(t, present, forced, "forced variable %q must always be present", forced)
			}
		})
	}
}

func TestPresentVariables_OptionalVarsFollowHeader(t *testing.T) {
	testCases := []struct {
		name           string
		header         []string
		expectVarName  string
		expectPresence bool
	}{
		{"body site present, exact case", []string{models.BodySiteVar.SourceColumnName}, models.BodySiteVar.Name, true},
		{"body site present, different case", []string{"SAMPLE ANATOMICAL LOCATION"}, models.BodySiteVar.Name, true},
		{"body site absent", []string{"unrelated"}, models.BodySiteVar.Name, false},
		{"laterality present", []string{models.LateralityVar.SourceColumnName}, models.LateralityVar.Name, true},
		{"laterality absent", []string{"unrelated"}, models.LateralityVar.Name, false},
		{"sample collection site present", []string{models.SampleCollectionSiteVar.SourceColumnName}, models.SampleCollectionSiteVar.Name, true},
		{"sample collection site absent", []string{"unrelated"}, models.SampleCollectionSiteVar.Name, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			present := namesOf(presentVariables(tc.header))
			if tc.expectPresence {
				assert.Contains(t, present, tc.expectVarName)
			} else {
				assert.NotContains(t, present, tc.expectVarName)
			}
		})
	}
}

func TestWriteFiles(t *testing.T) {
	samplesHeader := []string{samples.IDLabel, models.BodySiteVar.SourceColumnName, models.LateralityVar.SourceColumnName}
	consentedSamples := []samples.Sample{
		{ID: "sam-1", SubjectID: "sub-1", Values: map[string]string{models.BodySiteVar.SourceColumnName: "brain", models.LateralityVar.SourceColumnName: "left"}},
		{ID: "sam-2", SubjectID: "sub-2", Values: map[string]string{models.BodySiteVar.SourceColumnName: "liver", models.LateralityVar.SourceColumnName: "right"}},
	}

	outputDirectory := t.TempDir()
	require.NoError(t, WriteFiles(outputDirectory, config.RNA, config.YES, samplesHeader, consentedSamples))

	dsFile, err := excelize.OpenFile(outputDirectory + "/6a_SampleAttributes_DS.xlsx")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, dsFile.Close())
	}()

	records, err := dsFile.GetRows("Sheet1")
	require.NoError(t, err)
	require.Len(t, records, len(consentedSamples)+1)

	analyteTypeIdx := indexOfHeader(t, records[0], models.AnalyteTypeVar.Name)
	isTumorIdx := indexOfHeader(t, records[0], models.IsTumorVar.Name)
	for _, dataRow := range records[1:] {
		assert.Equal(t, config.RNA.String(), dataRow[analyteTypeIdx])
		assert.Equal(t, config.YES.String(), dataRow[isTumorIdx])
	}
}

func TestWriteFiles_ComputedColumnsPopulatedWithoutOptionalSourceColumns(t *testing.T) {
	samplesHeader := []string{samples.IDLabel}
	consentedSamples := []samples.Sample{
		{ID: "sam-1", SubjectID: "sub-1", Values: map[string]string{}},
	}

	outputDirectory := t.TempDir()
	require.NoError(t, WriteFiles(outputDirectory, config.DNARNA, config.NO, samplesHeader, consentedSamples))

	dsFile, err := excelize.OpenFile(outputDirectory + "/6a_SampleAttributes_DS.xlsx")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, dsFile.Close())
	}()

	records, err := dsFile.GetRows("Sheet1")
	require.NoError(t, err)
	require.Len(t, records, len(consentedSamples)+1)

	analyteTypeIdx := indexOfHeader(t, records[0], models.AnalyteTypeVar.Name)
	isTumorIdx := indexOfHeader(t, records[0], models.IsTumorVar.Name)
	assert.Equal(t, config.DNARNA.String(), records[1][analyteTypeIdx])
	assert.Equal(t, config.NO.String(), records[1][isTumorIdx])
}

func indexOfHeader(t *testing.T, header []string, name string) int {
	t.Helper()
	for i, h := range header {
		if h == name {
			return i
		}
	}
	require.Failf(t, "variable not found in header", "%q not found in %v", name, header)
	return -1
}
