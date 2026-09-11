package enumkit

import (
	"fmt"
	"strings"
)

// Enum backs an integer Enum's String()/FromString off one []string.
// The slice index IS the Enum's integer value, so names must be listed
// in iota order.
type Enum[T ~int] struct {
	typeName string
	names    []string
}

func NewEnum[T ~int](typeName string, names ...string) Enum[T] {
	return Enum[T]{typeName: typeName, names: names}
}

func (e Enum[T]) Name(v T) string {
	if i := int(v); i >= 0 && i < len(e.names) {
		return e.names[i]
	}
	return "UNKNOWN"
}

func (e Enum[T]) FromString(s string) (T, error) {
	for i, name := range e.names {
		if strings.EqualFold(s, name) {
			return T(i), nil
		}
	}
	return -1, fmt.Errorf("unknown %s '%s'", e.typeName, s)
}
