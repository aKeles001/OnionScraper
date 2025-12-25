package proxy

import (
	"net/http"

	"onionscraper/logic/config"

	"golang.org/x/net/proxy"
)

func NewTorClient(cfg config.Config) (*http.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", cfg.TorProxy, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}

	return client, nil
}
