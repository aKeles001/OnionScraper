package output

import (
	"os"
	"path/filepath"
)

type Writer struct {
	outputDir string
}

func NewWriter(outputDir string) (*Writer, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, err
	}
	return &Writer{outputDir: outputDir}, nil
}

func (w *Writer) WriteResult(target string, body []byte, screenshot []byte) error {
	// Sanitize target to create valid filenames
	safeTarget := filepath.Base(target)

	// Write HTML body
	htmlDir := filepath.Join(w.outputDir, "html")
	htmlPath := filepath.Join(htmlDir, safeTarget+".html")
	if err := os.WriteFile(htmlPath, body, 0644); err != nil {
		return err
	}

	// Write screenshot
	screenshotDir := filepath.Join(w.outputDir, "screenshots")
	screenshotPath := filepath.Join(screenshotDir, safeTarget+".png")
	if err := os.WriteFile(screenshotPath, screenshot, 0644); err != nil {
		return err
	}

	return nil
}
