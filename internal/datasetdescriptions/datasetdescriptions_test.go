package datasetdescriptions

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func TestFromFile(t *testing.T) {
	testCases := []struct {
		name           string
		filename       string
		expectedTitle  string
		expectedDOIURL string
	}{
		{"first version wording", "dataset_description__doi-for-first-version.xlsx", "Test Dataset Alpha", "https://doi.org/10.26275/bb00-frst"},
		{"this dataset wording", "dataset_description__doi-for-this-dataset.xlsx", "Test Dataset Alpha", "https://doi.org/10.26275/aa00-this"},
		{"first version with no 'the'", "dataset_description__doi-first-version-no-the.xlsx", "Test Dataset Gamma", "https://doi.org/10.26275/gg00-frst"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join("testdata", tc.filename)
			file, err := excelize.OpenFile(path)
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, file.Close())
			}()

			datasetDescription, err := FromFile(file)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedTitle, datasetDescription.Title)
			assert.Equal(t, tc.expectedDOIURL, datasetDescription.DOIURL)

		})
	}

}

func TestDOIURLIdentifierDescriptionsKeysAreMatchable(t *testing.T) {
	for k := range DOIURLIdentifierDescriptions {
		assert.Equal(t, k, strings.ToLower(k), "DOIURLIdentifierDescriptions key %q is not lowercase; lookups against it will never match", k)
		assert.Equal(t, k, strings.TrimSpace(k), "DOIURLIdentifierDescriptions key %q is not trimmed of whitespace; lookups against it will never match", k)
	}
}
