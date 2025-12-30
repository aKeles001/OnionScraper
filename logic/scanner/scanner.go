package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"onionscraper/logic/logger"
	"onionscraper/logic/output"
	"time"

	"github.com/chromedp/chromedp"
)

type Scanner struct {
	Client  *http.Client
	Writer  *output.Writer
	Timeout time.Duration
}
type Options struct {
	Targets []string
	Client  *http.Client
	Writer  *output.Writer
	Timeout time.Duration
	Retries int
}
type TorStatus struct {
	IP       string
	IsTor    bool
	Response string
}

func NewScanner(client *http.Client, writer *output.Writer, timeout time.Duration) *Scanner {
	return &Scanner{
		Client:  client,
		Writer:  writer,
		Timeout: timeout,
	}
}

func Run(opts Options) {
	scanner := NewScanner(opts.Client, opts.Writer, opts.Timeout)
	response, err := scanner.checkTorStatus()
	if err != nil {
		logger.Error("Failed to check Tor status: %v", err)
	} else {
		if response.IsTor {
			logger.Info("Connected to Tor network", "IP", response.IP)
		} else {
			logger.Error("Not connected to Tor network")
		}
	}
	for _, target := range opts.Targets {
		fmt.Printf("Scanning target: %s\n", target)
		for i := 0; i < opts.Retries; i++ {
			fmt.Printf("Scraping - Attempt %d/3\n", i+1)
			response, err := opts.Client.Get(target)
			if err != nil {
				logger.Error("Request failed", "error", err, "target", target, "attempt", i+1)
				time.Sleep(time.Duration(i+1) * 2 * time.Second)
				continue
			}
			if response.StatusCode == http.StatusOK {
				body, err := io.ReadAll(response.Body)
				if err != nil {
					logger.Error("Failed to read response body", "error", err, "target", target)
					response.Body.Close()
					continue
				}

				screenShot, err := scanner.CaptureScreenshot(target)
				if err != nil {
					logger.Error("Screenshot capture failed", "error", err, "target", target)
					response.Body.Close()
					continue
				}
				response.Body.Close()
				err = opts.Writer.WriteResult(target, body, screenShot)
				if err != nil {
					logger.Error("Failed to write result", "error", err, "target", target)
					continue
				}
				logger.Info("Successfully scraped target", "target", target)
				break
			}
		}
	}
}

func (s *Scanner) CaptureScreenshot(targetURL string) ([]byte, error) {
	if s.Timeout == 0 {
		s.Timeout = 90 * time.Second
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer("socks5://127.0.0.1:9050"),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.NoSandbox,
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, s.Timeout)
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),

		chromedp.WaitVisible(`body`, chromedp.ByQuery),

		chromedp.Sleep(25*time.Second),

		chromedp.FullScreenshot(&buf, 90),
	)

	if err != nil {
		return nil, fmt.Errorf("screenshot failed: %w", err)
	}
	return buf, nil
}

func (s *Scanner) checkTorStatus() (*TorStatus, error) {
	resp, err := s.Client.Get("https://check.torproject.org/api/ip")
	if err != nil {
		return nil, fmt.Errorf("Tor check failed: %w", err)
	}
	defer resp.Body.Close()
	var status TorStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}
	return &status, nil
}
