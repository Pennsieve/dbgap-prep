package ds

import (
	"encoding/csv"
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"log/slog"
	"os"
	"path/filepath"
)

type TSVWriter struct {
	path string
}

func NewTSVWriter(outputDirectory string, baseFileName string) *TSVWriter {
	return &TSVWriter{
		path: filepath.Join(outputDirectory, fmt.Sprintf("%s.txt", baseFileName)),
	}
}

func (t *TSVWriter) Path() string {
	return t.path
}

func (t *TSVWriter) Write(spec Spec, rows [][]string) error {
	file, err := os.Create(t.path)
	if err != nil {
		return fmt.Errorf("error creating DS file %s: %w", t.path, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			logger.Warn("error closing DS file", slog.String("path", t.path),
				slog.Any("error", err))
		}
	}()

	tsvWriter := csv.NewWriter(file)
	tsvWriter.Comma = '\t'
	// some users may use windows
	tsvWriter.UseCRLF = true

	records := make([][]string, 0, len(rows)+1)
	header := dd.VariableNames(spec.Variables)
	records = append(records, header)
	records = append(records, rows...)
	if err := tsvWriter.WriteAll(records); err != nil {
		return fmt.Errorf("error writing records to DS file %s: %s", t.path, err)
	}

	logger.Info("wrote DS file", slog.String("file", t.path))

	return nil
}
