package utils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/francomano/proxydoctor/core/check"
)

// ParseProxyConfig parses a proxy URL string (e.g. "http://user:pass@host:port")
// together with an explicit proxy type hint ("auto", "http", "https", "socks4", "socks5")
// into a check.ProxyConfig. An empty proxyStr means "no proxy" (direct connection).
//
// This is shared by the CLI (`cmd/cli/commands/diagnose.go`) and the HTTP server
// (`cmd/server/main.go`) so both entry points behave identically.
func ParseProxyConfig(proxyStr, proxyTypeStr string) (check.ProxyConfig, error) {
	if proxyStr == "" {
		return check.ProxyConfig{Type: check.ProxyTypeDirect}, nil
	}

	parsed, err := url.Parse(proxyStr)
	if err != nil {
		return check.ProxyConfig{}, fmt.Errorf("invalid proxy URL %q: %w", proxyStr, err)
	}
	if parsed.Host == "" {
		return check.ProxyConfig{}, fmt.Errorf("invalid proxy URL %q: missing host", proxyStr)
	}

	// Determine proxy type: explicit hint wins, otherwise infer from the URL scheme
	var proxyType check.ProxyType
	switch strings.ToLower(proxyTypeStr) {
	case "http":
		proxyType = check.ProxyTypeHTTP
	case "https":
		proxyType = check.ProxyTypeHTTPS
	case "socks4":
		proxyType = check.ProxyTypeSOCKS4
	case "socks5":
		proxyType = check.ProxyTypeSOCKS5
	case "auto", "":
		switch strings.ToLower(parsed.Scheme) {
		case "http":
			proxyType = check.ProxyTypeHTTP
		case "https":
			proxyType = check.ProxyTypeHTTPS
		case "socks4":
			proxyType = check.ProxyTypeSOCKS4
		case "socks5", "socks":
			proxyType = check.ProxyTypeSOCKS5
		default:
			return check.ProxyConfig{}, fmt.Errorf("cannot infer proxy type from scheme %q, pass an explicit proxy type", parsed.Scheme)
		}
	default:
		return check.ProxyConfig{}, fmt.Errorf("unknown proxy type %q", proxyTypeStr)
	}

	// Resolve host/port, applying sensible defaults per proxy type
	host := parsed.Hostname()
	portStr := parsed.Port()
	var port int
	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return check.ProxyConfig{}, fmt.Errorf("invalid proxy port %q: %w", portStr, err)
		}
		port = p
	} else {
		switch proxyType {
		case check.ProxyTypeHTTP:
			port = 8080
		case check.ProxyTypeHTTPS:
			port = 443
		case check.ProxyTypeSOCKS4, check.ProxyTypeSOCKS5:
			port = 1080
		}
	}

	config := check.ProxyConfig{
		Type: proxyType,
		Host: host,
		Port: port,
	}
	if parsed.User != nil {
		config.Username = parsed.User.Username()
		if pw, ok := parsed.User.Password(); ok {
			config.Password = pw
		}
	}

	return config, nil
}