package plugin

import (
	"fmt"
	"sync"

	"github.com/francomano/proxydoctor/core/engine"
)

// Plugin is the base interface every plugin must implement.
type Plugin interface {
	ID() string
	Name() string
	Version() string
	Description() string

	// Init is called once when the plugin is loaded.
	Init(ctx *Context) error

	// Shutdown is called when the plugin is unloaded.
	Shutdown() error
}

// Context is passed to Plugin.Init so plugins can access engine internals.
type Context struct {
	Registry *engine.CheckRegistry
	Config   map[string]interface{}
}

// CheckPlugin adds new diagnostic checks to the registry.
type CheckPlugin interface {
	Plugin
	RegisterChecks(registry *engine.CheckRegistry) error
}

// ExportPlugin adds new output formats (beyond json/text/markdown/html).
type ExportPlugin interface {
	Plugin
	FormatName() string
	FormatReport(report *engine.DiagnosisReport) (string, error)
}

// MiddlewarePlugin intercepts requests before and results after diagnosis.
type MiddlewarePlugin interface {
	Plugin
	BeforeDiagnosis(req *engine.DiagnosisRequest) error
	AfterDiagnosis(report *engine.DiagnosisReport) error
}

// Manager handles loading, registering, and lifecycle of plugins.
type Manager struct {
	plugins []Plugin
	mu      sync.RWMutex
}

// NewManager creates a new plugin manager.
func NewManager() *Manager {
	return &Manager{}
}

// Register adds a plugin to the manager and calls its Init.
func (m *Manager) Register(p Plugin, ctx *Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := p.Init(ctx); err != nil {
		return fmt.Errorf("plugin %s init failed: %w", p.ID(), err)
	}

	// Auto-register checks if the plugin implements CheckPlugin
	if cp, ok := p.(CheckPlugin); ok {
		if err := cp.RegisterChecks(ctx.Registry); err != nil {
			return fmt.Errorf("plugin %s register checks failed: %w", p.ID(), err)
		}
	}

	m.plugins = append(m.plugins, p)
	return nil
}

// Unregister removes a plugin and calls its Shutdown.
func (m *Manager) Unregister(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, p := range m.plugins {
		if p.ID() == id {
			if err := p.Shutdown(); err != nil {
				return fmt.Errorf("plugin %s shutdown failed: %w", p.ID(), err)
			}
			m.plugins = append(m.plugins[:i], m.plugins[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("plugin %s not found", id)
}

// Plugins returns all registered plugins.
func (m *Manager) Plugins() []Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Plugin, len(m.plugins))
	copy(result, m.plugins)
	return result
}

// GetPlugin returns a plugin by ID.
func (m *Manager) GetPlugin(id string) (Plugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, p := range m.plugins {
		if p.ID() == id {
			return p, true
		}
	}
	return nil, false
}

// ShutdownAll calls Shutdown on every registered plugin.
func (m *Manager) ShutdownAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for _, p := range m.plugins {
		if err := p.Shutdown(); err != nil {
			lastErr = fmt.Errorf("plugin %s shutdown failed: %w", p.ID(), err)
		}
	}
	return lastErr
}

// RunMiddlewareBefore runs BeforeDiagnosis on all middleware plugins.
func (m *Manager) RunMiddlewareBefore(req *engine.DiagnosisRequest) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, p := range m.plugins {
		if mw, ok := p.(MiddlewarePlugin); ok {
			if err := mw.BeforeDiagnosis(req); err != nil {
				return fmt.Errorf("middleware %s before failed: %w", p.ID(), err)
			}
		}
	}
	return nil
}

// RunMiddlewareAfter runs AfterDiagnosis on all middleware plugins.
func (m *Manager) RunMiddlewareAfter(report *engine.DiagnosisReport) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, p := range m.plugins {
		if mw, ok := p.(MiddlewarePlugin); ok {
			if err := mw.AfterDiagnosis(report); err != nil {
				return fmt.Errorf("middleware %s after failed: %w", p.ID(), err)
			}
		}
	}
	return nil
}

// GetExportFormats returns all export formats provided by export plugins.
func (m *Manager) GetExportFormats() map[string]ExportPlugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	formats := make(map[string]ExportPlugin)
	for _, p := range m.plugins {
		if ep, ok := p.(ExportPlugin); ok {
			formats[ep.FormatName()] = ep
		}
	}
	return formats
}
