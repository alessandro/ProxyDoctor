package engine

import (
	"sync"
	"time"

	"github.com/francomano/proxydoctor/core/check"
)

// DefaultExecutionContext is the default implementation of ExecutionContext
type DefaultExecutionContext struct {
	url            string
	proxyConfig    check.ProxyConfig
	directAdapter  check.NetworkAdapter
	proxyAdapter   check.NetworkAdapter
	sharedData     map[string]interface{}
	sharedDataLock sync.RWMutex
	timeout        time.Duration
	cancelled      bool
	cancelLock     sync.RWMutex
}

// NewExecutionContext creates a new execution context
func NewExecutionContext(
	url string,
	proxyConfig check.ProxyConfig,
	directAdapter check.NetworkAdapter,
	proxyAdapter check.NetworkAdapter,
	timeout time.Duration,
) check.ExecutionContext {
	return &DefaultExecutionContext{
		url:           url,
		proxyConfig:   proxyConfig,
		directAdapter: directAdapter,
		proxyAdapter:  proxyAdapter,
		sharedData:    make(map[string]interface{}),
		timeout:       timeout,
	}
}

func (c *DefaultExecutionContext) GetURL() string {
	return c.url
}

func (c *DefaultExecutionContext) GetProxyConfig() check.ProxyConfig {
	return c.proxyConfig
}

func (c *DefaultExecutionContext) GetDirectAdapter() check.NetworkAdapter {
	return c.directAdapter
}

func (c *DefaultExecutionContext) GetProxyAdapter() check.NetworkAdapter {
	return c.proxyAdapter
}

func (c *DefaultExecutionContext) GetSharedData(key string) interface{} {
	c.sharedDataLock.RLock()
	defer c.sharedDataLock.RUnlock()
	return c.sharedData[key]
}

func (c *DefaultExecutionContext) SetSharedData(key string, value interface{}) {
	c.sharedDataLock.Lock()
	defer c.sharedDataLock.Unlock()
	c.sharedData[key] = value
}

func (c *DefaultExecutionContext) GetTimeout() time.Duration {
	return c.timeout
}

func (c *DefaultExecutionContext) IsCancelled() bool {
	c.cancelLock.RLock()
	defer c.cancelLock.RUnlock()
	return c.cancelled
}

func (c *DefaultExecutionContext) Cancel() {
	c.cancelLock.Lock()
	defer c.cancelLock.Unlock()
	c.cancelled = true
}
