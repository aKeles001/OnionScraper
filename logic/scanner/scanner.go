package scanner

import (
	"net/http"
	"time"

	"onionscraper/logic/output"
)

type Scanner struct {
	Client  *http.Client
	Writer  *output.Writer
	Timeout time.Duration
}

func NewScanner(client *http.Client, writer *output.Writer, timeout time.Duration) *Scanner {
	return &Scanner{
		Client:  client,
		Writer:  writer,
		Timeout: timeout,
	}
}

type Options struct {
	Targets []string
	Client  *http.Client
	Writer  *output.Writer
	Timeout time.Duration
}

func Run(opts Options) {
	scanner := NewScanner(opts.Client, opts.Writer, opts.Timeout)

	for _, target := range opts.Targets {
		// Here would be the scanning logic for each target
		// For demonstration, we just create a dummy report
		reportData := []byte("Scan report for " + target)

		// Write the report using the Writer
		fileName := target + "_report.txt"
		if err := scanner.Writer.WriteReport(fileName, reportData); err != nil {
			// Handle error (e.g., log it)
			continue
		}

		// Optionally, add delays or other scanning logic here
		time.Sleep(1 * time.Second) // Dummy delay between scans
	}
}
