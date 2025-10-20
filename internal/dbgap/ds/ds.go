package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/logging"
)

var logger = logging.PackageLogger("ds")

type Spec struct {
	Variables []dd.Variable
}

type ToRowFunc[T any] func(variableNames []string, item T) []string

func ToRows[T any](variables []dd.Variable, items []T, toRow ToRowFunc[T]) [][]string {
	variableNames := dd.VariableNames(variables)

	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, toRow(variableNames, item))
	}
	return rows
}
