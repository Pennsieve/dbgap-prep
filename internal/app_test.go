package app

import (
	scds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/ds"
	scmodels "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/models"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetConsented(t *testing.T) {
	subjectConsents := []scds.SubjectConsent{
		{SubjectID: "sub-not-consented", Consent: scmodels.NoConsent.Value, Sex: scmodels.FemaleSex.Value},
		{SubjectID: "sub-1", Consent: scmodels.GRUConsent.Value, Sex: scmodels.MaleSex.Value},
		{SubjectID: "sub-3", Consent: scmodels.GRUConsent.Value, Sex: scmodels.UnknownSex.Value},
		{SubjectID: "sub-no-samps", Consent: scmodels.GRUConsent.Value, Sex: scmodels.FemaleSex.Value},
	}
	samps := []samples.Sample{
		{ID: "samp-1-sub-1", SubjectID: "sub-1"},
		{ID: "samp-1-sub-3", SubjectID: "sub-3"},
		{ID: "samp-2-sub-3", SubjectID: "sub-3"},
		{ID: "samp-3-sub-3", SubjectID: "sub-3"},
		{ID: "samp-1-no-sub"},
	}

	consentedSubjects, consentedSamples := GetConsented(subjectConsents, samps)

	// Subjects
	require.Len(t, consentedSubjects, 3)

	require.Contains(t, consentedSubjects, "sub-1")
	assert.Equal(t, "sub-1", consentedSubjects["sub-1"].SubjectID)

	require.Contains(t, consentedSubjects, "sub-3")
	assert.Equal(t, "sub-3", consentedSubjects["sub-3"].SubjectID)

	require.Contains(t, consentedSubjects, "sub-no-samps")
	assert.Equal(t, "sub-no-samps", consentedSubjects["sub-no-samps"].SubjectID)

	// Samples
	require.Len(t, consentedSamples, 2)

	require.Contains(t, consentedSamples, "sub-1")
	assert.Len(t, consentedSamples["sub-1"], 1)

	require.Contains(t, consentedSamples, "sub-3")
	assert.Len(t, consentedSamples["sub-3"], 3)

}
