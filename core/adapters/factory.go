package adapters

import (
	"github.com/francomano/proxydoctor/core/check"
)

// AdapterFactory creates network adapters
type AdapterFactory interface {
	CreateAdapter(proxyType check.ProxyType, config check.ProxyConfig) check.NetworkAdapter
}

// DefaultAdapterFactory is the default implementation
type DefaultAdapterFactory struct{}

// NewAdapterFactory creates a new adapter factory
func NewAdapterFactory() AdapterFactory {
	return &DefaultAdapterFactory{}
}

// CreateAdapter creates an adapter based on proxy type
func (f *DefaultAdapterFactory) CreateAdapter(proxyType check.ProxyType, config check.ProxyConfig) check.NetworkAdapter {
	switch proxyType {
	case check.ProxyTypeDirect:
		return NewDirectAdapter()
	case check.ProxyTypeHTTP:
		return NewHTTPProxyAdapter(config)
	case check.ProxyTypeHTTPS:
		return NewHTTPSProxyAdapter(config)
	case check.ProxyTypeSOCKS4:
		return NewSOCKS4Adapter(config)
	case check.ProxyTypeSOCKS5:
		return NewSOCKS5Adapter(config)
	default:
		return NewDirectAdapter()
	}
}
