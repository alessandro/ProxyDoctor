package checks

import (
	dnsresolve "github.com/francomano/proxydoctor/core/checks/dns_resolve"
	portscan "github.com/francomano/proxydoctor/core/checks/port_scan"
	publicip "github.com/francomano/proxydoctor/core/checks/public_ip"
	tlscert "github.com/francomano/proxydoctor/core/checks/tls_cert"
	"github.com/francomano/proxydoctor/core/engine"
)

// RegisterDefaults registers all built-in diagnostic checks into the provided registry.
func RegisterDefaults(registry *engine.CheckRegistry) {
	registry.Register(publicip.NewPublicIPCheck())
	registry.Register(dnsresolve.NewDNSResolveCheck())
	registry.Register(tlscert.NewTLSCertCheck())
	registry.Register(portscan.NewPortScanCheck())
}
