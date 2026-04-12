package httpx

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// NewClient 创建可选代理的 http client
func NewClient(proxyConnect string, timeoutSeconds int) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	client := &http.Client{Transport: transport}
	if timeoutSeconds > 0 {
		client.Timeout = time.Duration(timeoutSeconds) * time.Second
	}

	proxyConnect = strings.TrimSpace(proxyConnect)
	if proxyConnect == "" {
		return client
	}

	proxyURL, err := url.Parse(proxyConnect)
	if err != nil {
		return client
	}

	switch proxyURL.Scheme {
	case "http", "https":
		transport.Proxy = http.ProxyURL(proxyURL)
	case "socks5", "socks5h":
		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			proxyDialer, err := proxy.FromURL(proxyURL, dialer)
			if err != nil {
				return dialer.DialContext(ctx, network, addr)
			}
			return proxyDialer.Dial(network, addr)
		}
		transport.DialContext = dialContext
	}

	return client
}
