package analytetype

import (
	"github.com/pennsieve/dbgap-prep/internal/enums/internal/enumkit"
)

type Type int

const (
	DNA Type = iota
	RNA
	DNARNA
)

var types = enumkit.NewEnum[Type]("analyte type", "DNA", "RNA", "DNA/RNA")

func (at Type) String() string {
	return types.Name(at)
}

func FromString(s string) (Type, error) {
	return types.FromString(s)
}
