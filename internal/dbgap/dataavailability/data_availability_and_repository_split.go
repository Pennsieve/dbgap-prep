package dataavailability

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/pennsieve/dbgap-prep/internal/datasetdescriptions"
	"github.com/pennsieve/dbgap-prep/internal/logging"
)

var logger = logging.PackageLogger("dataavailability")

const DefaultFileNameBase = "Data_Availability_and_Repository_Split.txt"

const template = `Data Availability and Repository Split
 
This dbGaP study ([PHS_ACCESSION], [STUDY_TITLE]) contains the raw sequencing data.
 
All processed and derived data, detailed subject demographics, and additional descriptive metadata for these samples reside on the SPARC Portal(sparc.science), not in dbGaP. The SPARC record serves as the master metadata record for the collection and is updated as additional datasets are deposited and published.
 
To locate the processed data and full metadata for a given sample, see the SPARC_DATASET_DOI variable in this study's Sample Attributes, which records the DOI of the SPARC dataset that holds those data. The complete mapping of all subjects and samples in the collection is available on the SPARC Portal at [SPARC_DATASET_DOI_URL]. Sample identifiers in this study correspond directly to the sample id values in that dataset, and subject identifiers correspond to its subject id values.
`

func WriteFile(outputDirectory string, phsAccession string, description datasetdescriptions.DatasetDescription) error {
	path := filepath.Join(outputDirectory, DefaultFileNameBase)
	var replacements []string

	replacements = appendReplacement(replacements, "[PHS_ACCESSION]", phsAccession)
	replacements = appendReplacement(replacements, "[STUDY_TITLE]", description.Title)
	replacements = appendReplacement(replacements, "[SPARC_DATASET_DOI_URL]", description.DOIURL)

	outputFile := strings.NewReplacer(replacements...).Replace(template)
	// make the text file Windows friendly just in case.
	outputFile = windowsFriendlyLineEndings(outputFile)

	err := os.WriteFile(path, []byte(outputFile), 0644)

	if err != nil {
		return fmt.Errorf("error writing Data Availability and Repository Split file %s: %w", path, err)
	}
	logger.Info("wrote Data Availability and Repository Split file", slog.String("file", path))
	return nil
}

func appendReplacement(replacements []string, placeholder string, replacement string) []string {
	if trimmed := strings.TrimSpace(replacement); len(trimmed) == 0 {
		logger.Warn("no valid replacement provided; Data Availability and Repository Split will contain placeholder",
			slog.String("placeholder", placeholder),
			slog.String("rawReplacement", replacement),
		)
	} else {
		replacements = append(replacements, placeholder, trimmed)
	}
	return replacements
}

// windowsFriendlyLineEndings simply replaces all '\n's with '\r\n' in input and returns the result.
// It makes no effort to check if input already contains '\r\n'
// sequences, so it should only be used on strings that are '\n'-only.
func windowsFriendlyLineEndings(input string) string {
	return strings.ReplaceAll(input, "\n", "\r\n")
}
