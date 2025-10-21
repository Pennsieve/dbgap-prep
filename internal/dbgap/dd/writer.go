package dd

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/utils"
	"github.com/xuri/excelize/v2"
	"log/slog"
	"path/filepath"
)

type Writer interface {
	Path() string
	Write(spec Spec) error
}

type XLSXWriter struct {
	path string
}

func (x *XLSXWriter) Path() string {
	return x.path
}

func (x *XLSXWriter) Write(spec Spec) error {
	ddFile := excelize.NewFile()
	defer utils.CloseExcelFile(ddFile, logger)

	if err := ddFile.SetSheetName("Sheet1", spec.SheetName); err != nil {
		return fmt.Errorf("error setting %s sheet name: %w", spec.FileName, err)
	}
	if err := Populate(ddFile, spec.SheetName, spec); err != nil {
		return fmt.Errorf("error populating %s: %w", spec.FileName, err)
	}
	if err := ddFile.SaveAs(x.path); err != nil {
		return fmt.Errorf("error writing %s to %s: %w", spec.FileName, x.path, err)
	}
	logger.Info("wrote DD file", slog.String("file", x.path))

	return nil
}

func NewXLSXWriter(outputDirectory string, fileName string) *XLSXWriter {
	return &XLSXWriter{
		path: filepath.Join(outputDirectory, fileName),
	}
}

type NoOpWriter struct {
	path string
}

func (x *NoOpWriter) Path() string {
	return x.path
}

func (x *NoOpWriter) Write(spec Spec) error {
	logger.Info("no DD file written; DD output turned off", slog.String("path", x.path))
	return nil
}

func NewNoOpWriter(outputDirectory string, fileName string) *NoOpWriter {
	return &NoOpWriter{
		path: filepath.Join(outputDirectory, fileName),
	}
}
