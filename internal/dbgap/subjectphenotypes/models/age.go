package models

import (
	"fmt"
	"log/slog"
	"math"
	"regexp"
	"strconv"

	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

const AgeVariableName = "AGE"

var AgeVariable = dd.Variable{
	Name:             AgeVariableName,
	Description:      "Subject age in years (whole integer)",
	Type:             dd.IntegerType,
	SourceColumnName: "age",
}

func AgeFromSubject(subject subjects.Subject) *int {
	ageValue, found := subject.SearchValue(AgeVariableName)
	if !found {
		return nil
	}
	intAge, err := extractInt(ageValue)
	if err != nil {
		logger.Warn("error extracting integer age",
			slog.String("value", ageValue),
			slog.String("error", err.Error()))
		return nil
	}
	return &intAge
}

var numberRe = regexp.MustCompile(`-?\d*\.?\d+`)

func extractInt(s string) (int, error) {
	match := numberRe.FindString(s)
	if match == "" {
		return 0, fmt.Errorf("no number found in %q", s)
	}
	f, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing %q: %w", match, err)
	}
	return int(math.Floor(f)), nil
}
