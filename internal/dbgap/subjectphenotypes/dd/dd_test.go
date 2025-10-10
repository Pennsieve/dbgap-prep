package dd

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVariables(t *testing.T) {
	var attributeLabels = []string{"pool id", "subject experimental group", "age", "sex", "species", "strain", "RRID for strain", "age category", "also in dataset", "member of", "metadata only", "number of directly derived samples", "laboratory internal id", "date of birth", "age range (min)", "age range (max)", "body mass", "genotype", "phenotype", "handedness", "reference atlas", "experimental log file path", "experiment date", "disease or disorder", "intervention", "disease model", "protocol title", "protocol url or doi", "BMI", "Height", "Cause of death", "Ventillator use", "Ventillator time", "Time of death", "Date of death", "Hardy scale", "Cohort", "Smoking frequency", "Smoking duration", "Drinking status", "Smoke type", "Cocaine", "Alcohol (type)", "Alcohol (number)", "Alcohol (period)", "Alcohol (years)", "Alcohol (comments)", "Heroin use", "Heroin history", "Prescription abuse", "Signs of drug abuse", "Drugs for non-medical use", "Sex for money or drugs", "Seizures", "Illicit drug use", "Recreational drug history", "Major depression", "CMV total Ab", "EBV IgG Ab", "EBV IgM Ab", "HBcAB total", "HCV Ab", "HIV Ab", "ALS", "Alzheimer's or dementia", "Arthritis", "Cancer (5yr)", "Cancer (current)", "Non-metastatic cancer", "HIV", "Heart attack", "Ischemic heart disease", "Heart disease", "Diabetes type I", "Diabetes type II"}

	variables := Variables(attributeLabels)

	require.Len(t, variables, len(attributeLabels)+1)

	assert.Equal(t, *dd.SubjectIDVar, variables[0])
	assert.Equal(t, "X", variables[0].Attributes[dd.UniqueKeyColumn])

	for i, label := range attributeLabels {
		assert.Equal(t, label, variables[i+1].Name)

	}

}
