package ds

import (
	"encoding/csv"
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"log/slog"
	"os"
)

var logger = logging.PackageLogger("ds")

type Spec struct {
	FileName string
	Header   []string
}

func Write(path string, spec Spec, rows [][]string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating DS file %s: %w", path, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			logger.Warn("error closing DS file", slog.String("path", path),
				slog.Any("error", err))
		}
	}()

	tsvWriter := csv.NewWriter(file)
	tsvWriter.Comma = '\t'
	// some users may use windows
	tsvWriter.UseCRLF = true

	records := make([][]string, 0, len(rows)+1)
	records = append(records, spec.Header)
	records = append(records, rows...)
	if err := tsvWriter.WriteAll(records); err != nil {
		return fmt.Errorf("error writing records to DS file %s: %s", path, err)
	}

	return nil
}

type ToRowFunc[T any] func(variableNames []string, item T) []string

func ToRows[T any](variableNames []string, items []T, toRow ToRowFunc[T]) [][]string {
	rows := make([][]string, 0, len(items))
	for _, consentedSample := range items {
		rows = append(rows, toRow(variableNames, consentedSample))
	}
	return rows
}
