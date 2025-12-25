package output

import (
	"os"
	"path/filepath"
)

type Writer struct {
	OutputDir string
}

func NewWriter(outputDir string) (*Writer, error) {
	// Ensure the output directory exists
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, err
	}
	return &Writer{OutputDir: outputDir}, nil
}

func (w *Writer) WriteReport(fileName string, data []byte) error {
	filePath := filepath.Join(w.OutputDir, fileName)
	return os.WriteFile(filePath, data, 0644)
}
