package adapters

import (
	"fmt"
	"time"

	"github.com/francomano/proxydoctor/core/check"
)

// SOCKS4Adapter handles SOCKS4 proxy connections
type SOCKS4Adapter struct {
	config check.ProxyConfig
}

// NewSOCKS4Adapter creates a new SOCKS4 adapter
func NewSOCKS4Adapter(config check.ProxyConfig) check.NetworkAdapter {
	return &SOCKS4Adapter{config: config}
}

func (a *SOCKS4Adapter) Type() check.ProxyType {
	return check.ProxyTypeSOCKS4
}

func (a *SOCKS4Adapter) ExecuteHTTPRequest(req *check.HTTPRequest) (*check.HTTPResponse, error) {
	return nil, fmt.Errorf("SOCKS4 HTTP requests not yet implemented")
}

func (a *SOCKS4Adapter) FollowRedirects(url string, maxRedirects int) ([]check.RedirectStep, error) {
	return nil, fmt.Errorf("SOCKS4 redirect following not yet implemented")
}

func (a *SOCKS4Adapter) ResolveDNS(hostname string) ([]string, error) {
	return nil, fmt.Errorf("SOCKS4 DNS resolution not yet implemented")
}

func (a *SOCKS4Adapter) GetPublicIP() (string, error) {
	return "", fmt.Errorf("SOCKS4 public IP not yet implemented")
}

func (a *SOCKS4Adapter) TestPort(host string, port int, timeout time.Duration) (bool, error) {
	return false, fmt.Errorf("SOCKS4 port testing not yet implemented")
}

func (a *SOCKS4Adapter) GetTLSCertificate(hostname string) (*check.CertificateInfo, error) {
	return nil, fmt.Errorf("SOCKS4 TLS not yet implemented")
}

func (a *SOCKS4Adapter) GetTLSCipherSuite(hostname string) (string, error) {
	return "", fmt.Errorf("SOCKS4 TLS not yet implemented")
}

func (a *SOCKS4Adapter) GetTLSVersion(hostname string) (string, error) {
	return "", fmt.Errorf("SOCKS4 TLS not yet implemented")
}

// SOCKS5Adapter handles SOCKS5 proxy connections
type SOCKS5Adapter struct {
	config check.ProxyConfig
}

// NewSOCKS5Adapter creates a new SOCKS5 adapter
func NewSOCKS5Adapter(config check.ProxyConfig) check.NetworkAdapter {
	return &SOCKS5Adapter{config: config}
}

func (a *SOCKS5Adapter) Type() check.ProxyType {
	return check.ProxyTypeSOCKS5
}

func (a *SOCKS5Adapter) ExecuteHTTPRequest(req *check.HTTPRequest) (*check.HTTPResponse, error) {
	return nil, fmt.Errorf("SOCKS5 HTTP requests not yet implemented")
}

func (a *SOCKS5Adapter) FollowRedirects(url string, maxRedirects int) ([]check.RedirectStep, error) {
	return nil, fmt.Errorf("SOCKS5 redirect following not yet implemented")
}

func (a *SOCKS5Adapter) ResolveDNS(hostname string) ([]string, error) {
	return nil, fmt.Errorf("SOCKS5 DNS resolution not yet implemented")
}

func (a *SOCKS5Adapter) GetPublicIP() (string, error) {
	return "", fmt.Errorf("SOCKS5 public IP not yet implemented")
}

func (a *SOCKS5Adapter) TestPort(host string, port int, timeout time.Duration) (bool, error) {
	return false, fmt.Errorf("SOCKS5 port testing not yet implemented")
}

func (a *SOCKS5Adapter) GetTLSCertificate(hostname string) (*check.CertificateInfo, error) {
	return nil, fmt.Errorf("SOCKS5 TLS not yet implemented")
}

func (a *SOCKS5Adapter) GetTLSCipherSuite(hostname string) (string, error) {
	return "", fmt.Errorf("SOCKS5 TLS not yet implemented")
}

func (a *SOCKS5Adapter) GetTLSVersion(hostname string) (string, error) {
	return "", fmt.Errorf("SOCKS5 TLS not yet implemented")
}
