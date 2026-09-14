package dataavailability

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pennsieve/dbgap-prep/internal/datasetdescriptions"
)

const DefaultFileNameBase = "Data_Availability_and_Repository_Split.txt"

const template = `Data Availability and Repository Split
 
This dbGaP study ([PHS_ACCESSION], [STUDY_TITLE]) contains the raw sequencing data.
 
All processed and derived data, detailed subject demographics, and additional descriptive metadata for these samples reside on the SPARC Portal(sparc.science), not in dbGaP. The SPARC record serves as the master metadata record for the collection and is updated as additional datasets are deposited and published.
 
To locate the processed data and full metadata for a given sample, see the SPARC_DATASET_DOI variable in this study's Sample Attributes, which records the DOI of the SPARC dataset that holds those data. The complete mapping of all subjects and samples in the collection is available on the SPARC Portal at [SPARC_DATASET_DOI_URL]. Sample identifiers in this study correspond directly to the sample id values in that dataset, and subject identifiers correspond to its subject id values.
`

func WriteFile(outputDirectory string, phsAccession string, description datasetdescriptions.DatasetDescription) error {
	path := filepath.Join(outputDirectory, DefaultFileNameBase)
	var replacements []string
	if phsAccession = strings.TrimSpace(phsAccession); len(phsAccession) > 0 {
		replacements = append(replacements, "[PHS_ACCESSION]", phsAccession)
	}
	if title := strings.TrimSpace(description.Title); len(title) > 0 {
		replacements = append(replacements, "[STUDY_TITLE]", title)
	}
	if doiURL := strings.TrimSpace(description.DOIURL); len(doiURL) > 0 {
		replacements = append(replacements, "[SPARC_DATASET_DOI_URL]", doiURL)
	}

	outputFile := strings.NewReplacer(replacements...).Replace(template)

	err := os.WriteFile(path, []byte(outputFile), 0644)

	if err != nil {
		return fmt.Errorf("error writing Data Availability and Repository Split file %s: %w", path, err)
	}

	return nil
}
