package config

import "time"

type Config struct {
	// Input
	TargetFile string

	// Network
	TorProxy string
	Timeout  time.Duration

	// Output
	OutputDir  string
	ReportFile string

	// Execution
	Workers int
	Retries int
}
