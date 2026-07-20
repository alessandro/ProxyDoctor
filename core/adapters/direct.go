package adapters

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/francomano/proxydoctor/core/check"
)

// DirectAdapter handles direct network connections (no proxy)
type DirectAdapter struct {
	client *http.Client
}

// NewDirectAdapter creates a new direct connection adapter
func NewDirectAdapter() check.NetworkAdapter {
	return &DirectAdapter{
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (d *DirectAdapter) Type() check.ProxyType {
	return check.ProxyTypeDirect
}

func (d *DirectAdapter) ExecuteHTTPRequest(req *check.HTTPRequest) (*check.HTTPResponse, error) {
	httpReq, err := http.NewRequest(req.Method, req.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	startTime := time.Now()
	httpResp, err := d.client.Do(httpReq)
	duration := time.Since(startTime)

	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResp.Body.Close()

	// Read body
	buf := make([]byte, 1024*1024) // 1MB limit
	n, _ := httpResp.Body.Read(buf)

	return &check.HTTPResponse{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Header,
		Body:       buf[:n],
		Duration:   duration,
	}, nil
}

func (d *DirectAdapter) FollowRedirects(url string, maxRedirects int) ([]check.RedirectStep, error) {
	redirects := []check.RedirectStep{}
	currentURL := url

	for i := 0; i < maxRedirects; i++ {
		req, err := http.NewRequest("GET", currentURL, nil)
		if err != nil {
			return nil, err
		}

		client := &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			location := resp.Header.Get("Location")
			if location == "" {
				break
			}

			redirects = append(redirects, check.RedirectStep{
				From:       currentURL,
				To:         location,
				StatusCode: resp.StatusCode,
				Headers:    resp.Header,
			})

			currentURL = location
		} else {
			break
		}
	}

	return redirects, nil
}

func (d *DirectAdapter) ResolveDNS(hostname string) ([]string, error) {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return nil, fmt.Errorf("DNS lookup failed: %w", err)
	}

	result := make([]string, len(ips))
	for i, ip := range ips {
		result[i] = ip.String()
	}
	return result, nil
}

func (d *DirectAdapter) GetPublicIP() (string, error) {
	// This would call an external service like ipify.org
	// For now, return a placeholder
	return "", fmt.Errorf("public IP retrieval not implemented for direct adapter")
}

func (d *DirectAdapter) TestPort(host string, port int, timeout time.Duration) (bool, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false, nil
	}
	defer conn.Close()
	return true, nil
}

func (d *DirectAdapter) GetTLSCertificate(hostname string) (*check.CertificateInfo, error) {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", hostname+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, fmt.Errorf("TLS handshake failed: %w", err)
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("no certificates received")
	}

	cert := certs[0]
	return &check.CertificateInfo{
		Subject:            cert.Subject.String(),
		Issuer:             cert.Issuer.String(),
		NotBefore:          cert.NotBefore,
		NotAfter:           cert.NotAfter,
		SANs:               cert.DNSNames,
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		PublicKeyType:      fmt.Sprintf("%T", cert.PublicKey),
		IsValid:            time.Now().Before(cert.NotAfter) && time.Now().After(cert.NotBefore),
	}, nil
}

func (d *DirectAdapter) GetTLSCipherSuite(hostname string) (string, error) {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", hostname+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return "", fmt.Errorf("TLS handshake failed: %w", err)
	}
	defer conn.Close()

	state := conn.ConnectionState()
	return tls.CipherSuiteName(state.CipherSuite), nil
}

func (d *DirectAdapter) GetTLSVersion(hostname string) (string, error) {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", hostname+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return "", fmt.Errorf("TLS handshake failed: %w", err)
	}
	defer conn.Close()

	state := conn.ConnectionState()
	return fmt.Sprintf("TLS %d.%d", state.Version>>8, state.Version&0xFF), nil
}
