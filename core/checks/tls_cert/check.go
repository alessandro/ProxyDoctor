package tlscert

import (
	"fmt"
	"net/url"
	"time"

	"github.com/francomano/proxydoctor/core/check"
)

type TLSCertCheck struct{}

func NewTLSCertCheck() check.Checker {
	return &TLSCertCheck{}
}

func (c *TLSCertCheck) ID() string         { return "tls_certificate" }
func (c *TLSCertCheck) Name() string        { return "TLS Certificate" }
func (c *TLSCertCheck) Description() string { return "Checks the TLS certificate of the target (validity, issuer, expiry, cipher suite, TLS version)" }
func (c *TLSCertCheck) Category() check.CheckCategory { return check.CategoryTLS }
func (c *TLSCertCheck) DependsOn() []string           { return []string{} }

func (c *TLSCertCheck) Execute(ctx check.ExecutionContext) check.CheckResult {
	result := check.NewCheckResult(c.ID(), c.Category())
	startTime := time.Now()

	targetURL := ctx.GetURL()
	parsed, err := url.Parse(targetURL)
	if err != nil {
		result.SetExecutionTime(time.Since(startTime))
		return *result.WithStatus(check.StatusError, check.SeverityCritical).
			WithExplanation(fmt.Sprintf("Invalid URL: %v", err)).
			WithConfidence(0)
	}

	hostname := parsed.Hostname()

	adapter := ctx.GetProxyAdapter()
	if ctx.GetProxyConfig().Type == check.ProxyTypeDirect {
		adapter = ctx.GetDirectAdapter()
	}

	// Only check TLS for HTTPS URLs
	if parsed.Scheme != "https" {
		result.SetExecutionTime(time.Since(startTime))
		return *result.WithStatus(check.StatusSkipped, check.SeverityInfo).
			WithExplanation(fmt.Sprintf("Skipping TLS check for non-HTTPS URL: %s", targetURL)).
			WithConfidence(0)
	}

	certInfo, err := adapter.GetTLSCertificate(hostname)
	if err != nil {
		result.SetExecutionTime(time.Since(startTime))
		return *result.WithStatus(check.StatusError, check.SeverityWarning).
			WithExplanation(fmt.Sprintf("TLS certificate check failed for %s: %v", hostname, err)).
			WithConfidence(0).
			AddProbableCause("Target may not support TLS").
			AddProbableCause("Proxy may be blocking TLS connections")
	}

	cipherSuite, _ := adapter.GetTLSCipherSuite(hostname)
	tlsVersion, _ := adapter.GetTLSVersion(hostname)

	ctx.SetSharedData("tls_cert", certInfo)
	ctx.SetSharedData("tls_cipher", cipherSuite)
	ctx.SetSharedData("tls_version", tlsVersion)

	if !certInfo.IsValid {
		result.SetExecutionTime(time.Since(startTime))
		return *result.WithStatus(check.StatusFailed, check.SeverityCritical).
			WithExplanation(fmt.Sprintf("TLS certificate for %s is invalid — not valid or expired", hostname)).
			WithConfidence(0.95).
			AddEvidence("subject", certInfo.Subject).
			AddEvidence("issuer", certInfo.Issuer).
			AddEvidence("not_before", certInfo.NotBefore).
			AddEvidence("not_after", certInfo.NotAfter).
			AddEvidence("is_valid", certInfo.IsValid).
			AddEvidence("cipher_suite", cipherSuite).
			AddEvidence("tls_version", tlsVersion).
			AddSuggestedAction("Check if the certificate has been revoked or if the system clock is correct")
	}

	result.SetExecutionTime(time.Since(startTime))
	return *result.WithStatus(check.StatusPassed, check.SeverityInfo).
		WithExplanation(fmt.Sprintf("TLS certificate for %s is valid (issuer: %s, expires: %s)", hostname, certInfo.Issuer, certInfo.NotAfter.Format("2006-01-02"))).
		WithConfidence(0.95).
		AddEvidence("subject", certInfo.Subject).
		AddEvidence("issuer", certInfo.Issuer).
		AddEvidence("not_before", certInfo.NotBefore).
		AddEvidence("not_after", certInfo.NotAfter).
		AddEvidence("is_valid", certInfo.IsValid).
		AddEvidence("signature_algorithm", certInfo.SignatureAlgorithm).
		AddEvidence("sans", certInfo.SANs).
		AddEvidence("cipher_suite", cipherSuite).
		AddEvidence("tls_version", tlsVersion)
}
