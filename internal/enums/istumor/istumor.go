package istumor

import (
	"github.com/pennsieve/dbgap-prep/internal/enums/internal/enumkit"
)

type Value int

const (
	NO Value = iota
	YES
)

var values = enumkit.NewEnum[Value]("is tumor", "No", "Yes")

func (it Value) String() string {
	return values.Name(it)
}

func FromString(s string) (Value, error) {
	return values.FromString(s)
}
