package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/logging"
)

var logger = logging.PackageLogger("ds")

type Spec struct {
	Variables []dd.Variable
}

type ToRowFunc[T any] func(variables []dd.Variable, item T) []string

func ToRows[T any](variables []dd.Variable, items []T, toRow ToRowFunc[T]) [][]string {

	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, toRow(variables, item))
	}
	return rows
}
