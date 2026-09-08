package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type FlexBool bool

func (b *FlexBool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var realBool bool
	err := json.Unmarshal(data, &realBool)
	if err == nil {
		*b = FlexBool(realBool)
		return nil
	}
	// A type mismatch is the expected "not a JSON bool" case — fall through and
	// try the stringified form. Anything else (e.g. a syntax error from a direct
	// call with malformed bytes) is genuinely unexpected, so surface it.
	var typeErr *json.UnmarshalTypeError
	if !errors.As(err, &typeErr) {
		return fmt.Errorf("FlexBool: unexpected error decoding bool: %w", err)
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("FlexBool: expected bool or string, got %s", data)
	}
	parsed, err := strconv.ParseBool(strings.TrimSpace(s))
	if err != nil {
		return fmt.Errorf("FlexBool: invalid bool string %q: %w", s, err)
	}
	*b = FlexBool(parsed)
	return nil
}
