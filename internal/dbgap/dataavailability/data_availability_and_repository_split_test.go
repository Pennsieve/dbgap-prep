package dataavailability

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pennsieve/dbgap-prep/internal/datasetdescriptions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteFile(t *testing.T) {
	outputDir := t.TempDir()

	phsAccession := "phs003456.v1.p1"
	description := datasetdescriptions.DatasetDescription{
		Title:  "A Study Of Things",
		DOIURL: "https://doi.org/10.26275/abcd-efgh",
	}

	require.NoError(t, WriteFile(outputDir, phsAccession, description))

	actual := readOutputFile(t, outputDir)

	assert.Contains(t, actual, phsAccession)
	assert.Contains(t, actual, description.Title)
	assert.Contains(t, actual, description.DOIURL)
	assertNoPlaceholdersRemain(t, actual)
}

// TestWriteFile_Empty checks that empty or blank values leave the template
// untouched instead of substituting empty strings or panicking.
func TestWriteFile_Empty(t *testing.T) {
	for name, testCase := range map[string]struct {
		phsAccession string
		description  datasetdescriptions.DatasetDescription
	}{
		"empty values": {},
		"blank values": {
			phsAccession: "  ",
			description:  datasetdescriptions.DatasetDescription{Title: "\t", DOIURL: " \n "},
		},
	} {
		t.Run(name, func(t *testing.T) {
			outputDir := t.TempDir()

			require.NoError(t, WriteFile(outputDir, testCase.phsAccession, testCase.description))

			assert.Equal(t, template, readOutputFile(t, outputDir))
		})
	}
}

func TestWriteFile_PartialValues(t *testing.T) {
	outputDir := t.TempDir()

	phsAccession := "phs003456.v1.p1"

	require.NoError(t, WriteFile(outputDir, phsAccession, datasetdescriptions.DatasetDescription{}))

	actual := readOutputFile(t, outputDir)

	assert.Contains(t, actual, phsAccession)
	assert.NotContains(t, actual, "[PHS_ACCESSION]")
	assert.Contains(t, actual, "[STUDY_TITLE]")
	assert.Contains(t, actual, "[SPARC_DATASET_DOI_URL]")
}

func TestWriteFile_TrimsValues(t *testing.T) {
	outputDir := t.TempDir()

	description := datasetdescriptions.DatasetDescription{
		Title:  "  A Study Of Things  ",
		DOIURL: "\thttps://doi.org/10.26275/abcd-efgh\n",
	}

	require.NoError(t, WriteFile(outputDir, "  phs003456.v1.p1\n", description))

	actual := readOutputFile(t, outputDir)

	assert.Contains(t, actual, "(phs003456.v1.p1, A Study Of Things)")
	assert.Contains(t, actual, "at https://doi.org/10.26275/abcd-efgh.")
	assertNoPlaceholdersRemain(t, actual)
}

func TestWriteFile_NonExistentOutputDirectory(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "does-not-exist")

	assert.Error(t, WriteFile(outputDir, "phs003456.v1.p1", datasetdescriptions.DatasetDescription{}))
}

func readOutputFile(t *testing.T, outputDir string) string {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(outputDir, DefaultFileNameBase))
	require.NoError(t, err)
	return string(contents)
}

func assertNoPlaceholdersRemain(t *testing.T, actual string) {
	t.Helper()
	for _, placeholder := range []string{"[PHS_ACCESSION]", "[STUDY_TITLE]", "[SPARC_DATASET_DOI_URL]"} {
		assert.NotContains(t, actual, placeholder)
	}
}
