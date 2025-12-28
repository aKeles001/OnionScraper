package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"onionscraper/logic/config"
	"onionscraper/logic/input"
	"onionscraper/logic/logger"
	"onionscraper/logic/output"
	"onionscraper/logic/proxy"
	"onionscraper/logic/scanner"
)

func main() {
	/*Command-line flags*/
	targetFile := flag.String("targets", "data/targets.yaml", "Path to targets file")
	outputDir := flag.String("out", "data/outputs", "Output directory")
	timeout := flag.Int("timeout", 30, "Request timeout in seconds")
	torProxy := flag.String("proxy", "127.0.0.1:9050", "Tor SOCKS5 proxy address")
	maxRetries := flag.Int("retries", 3, "Maximum number of retries for requests")
	help := flag.Bool("help", false, "Show help message")

	if *help {
		flag.Usage()
		return
	}
	flag.Parse()

	/*Load Configuration*/
	cfg := config.Config{
		TargetFile: *targetFile,
		OutputDir:  *outputDir,
		TorProxy:   *torProxy,
		Timeout:    time.Duration(*timeout) * time.Second,
		MaxRetries: *maxRetries,
	}

	if err := logger.Init(cfg.OutputDir); err != nil {
		fmt.Fprintf(os.Stderr, "logger init failed: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	/*Read Targets*/
	targets, err := input.ReadTargets(cfg.TargetFile)
	if err != nil {
		logger.Error("Failed to read targets: %v", err)
		os.Exit(1)
	}

	if len(targets) == 0 {
		logger.Error("No targets found in the specified file")
		os.Exit(1)
	}
	logger.Info("Loaded targets", "count", len(targets))
	for _, target := range targets {
		fmt.Println("Target", "url", target)
	}

	logger.Info("Tor Scraper started")

	client, err := proxy.TorClient(cfg)
	if err != nil {
		logger.Error("Failed to initialize Tor client: %v", err)
		os.Exit(1)
	}

	writer, err := output.NewWriter(cfg.OutputDir)
	if err != nil {
		logger.Error("Failed to initialize output writer: %v", err)
		os.Exit(1)
	}

	scanner.Run(scanner.Options{
		Targets: targets,
		Client:  client,
		Writer:  writer,
		Timeout: cfg.Timeout * 25,
		Retries: cfg.MaxRetries,
	})

	logger.Info("Scan completed")
}
