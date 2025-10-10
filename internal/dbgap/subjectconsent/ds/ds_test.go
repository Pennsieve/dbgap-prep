package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/models"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetConsented(t *testing.T) {
	subjectConsents := []SubjectConsent{
		{SubjectID: "sub-not-consented", Consent: models.NoConsent.Value, Sex: models.FemaleSex.Value},
		{SubjectID: "sub-1", Consent: models.GRUConsent.Value, Sex: models.MaleSex.Value},
		{SubjectID: "sub-3", Consent: models.GRUConsent.Value, Sex: models.UnknownSex.Value},
		{SubjectID: "sub-no-samps", Consent: models.GRUConsent.Value, Sex: models.FemaleSex.Value},
	}
	samps := []samples.Sample{
		{ID: "samp-2-sub-3", SubjectID: "sub-3"},
		{ID: "samp-1-sub-1", SubjectID: "sub-1"},
		{ID: "samp-1-no-sub"},
		{ID: "samp-1-sub-3", SubjectID: "sub-3"},
		{ID: "samp-1-not-consented", SubjectID: "sub-not-consented"},
		{ID: "samp-3-sub-3", SubjectID: "sub-3"},
	}

	consentedSubjects, consentedSamples, consentedSamplesInSubjectOrder := GetConsented(subjectConsents, samps)

	// Subjects
	require.Len(t, consentedSubjects, 3)

	assert.Equal(t, "sub-1", consentedSubjects[0].SubjectID)

	assert.Equal(t, "sub-3", consentedSubjects[1].SubjectID)

	assert.Equal(t, "sub-no-samps", consentedSubjects[2].SubjectID)

	// Samples
	require.Len(t, consentedSamples, 4)

	assert.Equal(t, "samp-2-sub-3", consentedSamples[0].ID)
	assert.Equal(t, "samp-1-sub-1", consentedSamples[1].ID)
	assert.Equal(t, "samp-1-sub-3", consentedSamples[2].ID)
	assert.Equal(t, "samp-3-sub-3", consentedSamples[3].ID)

	// Samples in subject order should have same elements as consentedSamples, but
	// be in the order determined by their subject id's subjectConsents position.
	assert.ElementsMatch(t, consentedSamples, consentedSamplesInSubjectOrder)
	assert.Equal(t, "samp-1-sub-1", consentedSamplesInSubjectOrder[0].ID)
	assert.Equal(t, "samp-2-sub-3", consentedSamplesInSubjectOrder[1].ID)
	assert.Equal(t, "samp-1-sub-3", consentedSamplesInSubjectOrder[2].ID)
	assert.Equal(t, "samp-3-sub-3", consentedSamplesInSubjectOrder[3].ID)

}
