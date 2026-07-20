package adapters

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/francomano/proxydoctor/core/check"
)

// HTTPProxyAdapter handles HTTP proxy connections
type HTTPProxyAdapter struct {
	config check.ProxyConfig
	client *http.Client
}

// NewHTTPProxyAdapter creates a new HTTP proxy adapter
func NewHTTPProxyAdapter(config check.ProxyConfig) check.NetworkAdapter {
	proxyURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
	parsedURL, _ := url.Parse(proxyURL)

	transport := &http.Transport{
		Proxy: http.ProxyURL(parsedURL),
	}

	return &HTTPProxyAdapter{
		config: config,
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (a *HTTPProxyAdapter) Type() check.ProxyType {
	return check.ProxyTypeHTTP
}

func (a *HTTPProxyAdapter) ExecuteHTTPRequest(req *check.HTTPRequest) (*check.HTTPResponse, error) {
	httpReq, err := http.NewRequest(req.Method, req.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	startTime := time.Now()
	httpResp, err := a.client.Do(httpReq)
	duration := time.Since(startTime)

	if err != nil {
		return nil, fmt.Errorf("HTTP request via proxy failed: %w", err)
	}
	defer httpResp.Body.Close()

	buf := make([]byte, 1024*1024)
	n, _ := httpResp.Body.Read(buf)

	return &check.HTTPResponse{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Header,
		Body:       buf[:n],
		Duration:   duration,
	}, nil
}

func (a *HTTPProxyAdapter) FollowRedirects(url string, maxRedirects int) ([]check.RedirectStep, error) {
	// Implementation similar to DirectAdapter
	return []check.RedirectStep{}, nil
}

func (a *HTTPProxyAdapter) ResolveDNS(hostname string) ([]string, error) {
	// HTTP proxy doesn't expose DNS directly
	return nil, fmt.Errorf("DNS resolution via HTTP proxy not directly available")
}

func (a *HTTPProxyAdapter) GetPublicIP() (string, error) {
	return "", fmt.Errorf("public IP retrieval not implemented for HTTP proxy")
}

func (a *HTTPProxyAdapter) TestPort(host string, port int, timeout time.Duration) (bool, error) {
	return false, fmt.Errorf("port testing via HTTP proxy not implemented")
}

func (a *HTTPProxyAdapter) GetTLSCertificate(hostname string) (*check.CertificateInfo, error) {
	return nil, fmt.Errorf("TLS certificate retrieval via HTTP proxy not implemented")
}

func (a *HTTPProxyAdapter) GetTLSCipherSuite(hostname string) (string, error) {
	return "", fmt.Errorf("TLS cipher suite detection via HTTP proxy not implemented")
}

func (a *HTTPProxyAdapter) GetTLSVersion(hostname string) (string, error) {
	return "", fmt.Errorf("TLS version detection via HTTP proxy not implemented")
}

// HTTPSProxyAdapter handles HTTPS proxy connections
type HTTPSProxyAdapter struct {
	config check.ProxyConfig
	client *http.Client
}

// NewHTTPSProxyAdapter creates a new HTTPS proxy adapter
func NewHTTPSProxyAdapter(config check.ProxyConfig) check.NetworkAdapter {
	proxyURL := fmt.Sprintf("https://%s:%d", config.Host, config.Port)
	parsedURL, _ := url.Parse(proxyURL)

	transport := &http.Transport{
		Proxy: http.ProxyURL(parsedURL),
	}

	return &HTTPSProxyAdapter{
		config: config,
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
}

func (a *HTTPSProxyAdapter) Type() check.ProxyType {
	return check.ProxyTypeHTTPS
}

func (a *HTTPSProxyAdapter) ExecuteHTTPRequest(req *check.HTTPRequest) (*check.HTTPResponse, error) {
	httpReq, err := http.NewRequest(req.Method, req.URL, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	startTime := time.Now()
	httpResp, err := a.client.Do(httpReq)
	duration := time.Since(startTime)

	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	buf := make([]byte, 1024*1024)
	n, _ := httpResp.Body.Read(buf)

	return &check.HTTPResponse{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Header,
		Body:       buf[:n],
		Duration:   duration,
	}, nil
}

func (a *HTTPSProxyAdapter) FollowRedirects(url string, maxRedirects int) ([]check.RedirectStep, error) {
	return []check.RedirectStep{}, nil
}

func (a *HTTPSProxyAdapter) ResolveDNS(hostname string) ([]string, error) {
	return nil, fmt.Errorf("DNS resolution via HTTPS proxy not directly available")
}

func (a *HTTPSProxyAdapter) GetPublicIP() (string, error) {
	return "", fmt.Errorf("public IP retrieval not implemented for HTTPS proxy")
}

func (a *HTTPSProxyAdapter) TestPort(host string, port int, timeout time.Duration) (bool, error) {
	return false, fmt.Errorf("port testing via HTTPS proxy not implemented")
}

func (a *HTTPSProxyAdapter) GetTLSCertificate(hostname string) (*check.CertificateInfo, error) {
	return nil, fmt.Errorf("TLS certificate retrieval via HTTPS proxy not implemented")
}

func (a *HTTPSProxyAdapter) GetTLSCipherSuite(hostname string) (string, error) {
	return "", fmt.Errorf("TLS cipher suite detection via HTTPS proxy not implemented")
}

func (a *HTTPSProxyAdapter) GetTLSVersion(hostname string) (string, error) {
	return "", fmt.Errorf("TLS version detection via HTTPS proxy not implemented")
}
