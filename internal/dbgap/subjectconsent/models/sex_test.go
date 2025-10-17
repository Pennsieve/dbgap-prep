package models

import (
	"github.com/google/uuid"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSexFromSubject(t *testing.T) {
	tests := []struct {
		scenario          string
		inputSex          string
		expectedOutputSex string
	}{
		{"empty", "", UnknownSex.Value},
		{"one", "Male", MaleSex.Value},
		{"two", "female", FemaleSex.Value},
		{"three", "3", UnknownSex.Value},
		{"four", "4", UnknownSex.Value},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			subject := subjects.Subject{
				ID:  uuid.NewString(),
				Sex: tt.inputSex,
			}
			outputSex := SexFromSubject(subject)
			assert.Equal(t, tt.expectedOutputSex, outputSex)
		})
	}
}
