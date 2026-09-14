package lambda

import (
	"testing"

	"github.com/pennsieve/dbgap-prep/internal/enums/analytetype"
	"github.com/pennsieve/dbgap-prep/internal/enums/consentgroup"
	"github.com/pennsieve/dbgap-prep/internal/enums/istumor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validEvent returns an Event that configFromEvent accepts. The enum values differ from
// their zero values so that an assertion on one of them fails if the field is never
// assigned.
func validEvent() Event {
	return Event{
		IntegrationID:      "integration-1",
		WorkflowInstanceID: "workflow-instance-1",
		InputDirectory:     "/tmp/input",
		OutputDirectory:    "/tmp/output",
		ConsentGroup:       consentgroup.HMB.String(),
		AnalyteType:        analytetype.RNA.String(),
		IsTumor:            istumor.YES.String(),
		PHSAccession:       "phs003456.v1.p1",
	}
}

func TestConfigFromEvent(t *testing.T) {
	event := validEvent()

	cfg, err := configFromEvent(event)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "integration-1", cfg.IntegrationID)
	assert.Equal(t, "workflow-instance-1", cfg.WorkflowInstanceID)
	assert.Equal(t, "/tmp/input", cfg.InputDirectory)
	assert.Equal(t, "/tmp/output", cfg.OutputDirectory)
	assert.Equal(t, consentgroup.HMB, cfg.ConsentGroup)
	assert.Equal(t, analytetype.RNA, cfg.AnalyteType)
	assert.Equal(t, istumor.YES, cfg.IsTumor)
	assert.Equal(t, "phs003456.v1.p1", cfg.PHSAccession)
}

func TestConfigFromEvent_EmptyPHSAccession(t *testing.T) {
	event := validEvent()
	event.PHSAccession = ""

	cfg, err := configFromEvent(event)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Empty(t, cfg.PHSAccession)
}

func TestConfigFromEvent_IsTumorOptional(t *testing.T) {
	for name, isTumor := range map[string]string{
		"empty":      "",
		"whitespace": "  ",
	} {
		t.Run(name, func(t *testing.T) {
			event := validEvent()
			// the default can only be the source of a NO below if the event starts as YES
			require.Equal(t, istumor.YES.String(), event.IsTumor)
			event.IsTumor = isTumor

			cfg, err := configFromEvent(event)
			require.NoError(t, err)
			require.NotNil(t, cfg)

			assert.Equal(t, istumor.NO, cfg.IsTumor)
		})
	}
}

func TestConfigFromEvent_IsTumorTrimmed(t *testing.T) {
	event := validEvent()
	event.IsTumor = "  " + istumor.NO.String() + "  "

	cfg, err := configFromEvent(event)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, istumor.NO, cfg.IsTumor)
}

// TestConfigFromEvent_Errors covers what configFromEvent itself is responsible for:
// converting each event field and passing the failure on. The parsing of individual
// values is covered by the tests of the enum packages.
func TestConfigFromEvent_Errors(t *testing.T) {
	for name, modify := range map[string]func(event *Event){
		"missing consent group": func(event *Event) { event.ConsentGroup = "" },
		"invalid consent group": func(event *Event) { event.ConsentGroup = "not-a-consent-group" },
		"missing analyte type":  func(event *Event) { event.AnalyteType = "" },
		"invalid analyte type":  func(event *Event) { event.AnalyteType = "not-an-analyte-type" },
		"invalid is tumor":      func(event *Event) { event.IsTumor = "not-an-is-tumor" },
	} {
		t.Run(name, func(t *testing.T) {
			event := validEvent()
			modify(&event)

			cfg, err := configFromEvent(event)

			assert.Error(t, err)
			assert.Nil(t, cfg)
		})
	}
}
