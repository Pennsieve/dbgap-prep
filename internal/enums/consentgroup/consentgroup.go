package consentgroup

import (
	"github.com/pennsieve/dbgap-prep/internal/enums/internal/enumkit"
)

type Group int

const (
	GRU Group = iota
	HMB
	OTHER
)

// groups order must match enum order listed above.
var groups = enumkit.NewEnum[Group]("consent group", "GRU", "HMB", "Other")

func (g Group) String() string {
	return groups.Name(g)
}

func FromString(s string) (Group, error) {
	return groups.FromString(s)
}
