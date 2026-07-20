package engine

import (
	"fmt"
	"sync"

	"github.com/francomano/proxydoctor/core/check"
)

// DependencyDAG represents a directed acyclic graph of check dependencies
type DependencyDAG interface {
	AddNode(id string, checker check.Checker)
	AddEdge(from, to string)
	TopologicalSort() []string
	HasCycle() bool
}

// DefaultDAG is the default implementation of DependencyDAG
type DefaultDAG struct {
	nodes map[string]check.Checker
	edges map[string][]string
	mu    sync.RWMutex
}

// NewDependencyDAG creates a new dependency DAG
func NewDependencyDAG() DependencyDAG {
	return &DefaultDAG{
		nodes: make(map[string]check.Checker),
		edges: make(map[string][]string),
	}
}

func (d *DefaultDAG) AddNode(id string, checker check.Checker) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.nodes[id] = checker
	if _, exists := d.edges[id]; !exists {
		d.edges[id] = []string{}
	}
}

func (d *DefaultDAG) AddEdge(from, to string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.edges[from]; !exists {
		d.edges[from] = []string{}
	}
	d.edges[from] = append(d.edges[from], to)

	if _, exists := d.edges[to]; !exists {
		d.edges[to] = []string{}
	}
}

// TopologicalSort returns the checks in topological order (dependencies first)
func (d *DefaultDAG) TopologicalSort() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	visited := make(map[string]bool)
	stack := []string{}

	var dfs func(id string)
	dfs = func(id string) {
		if visited[id] {
			return
		}
		visited[id] = true

		// Visit all nodes that depend on this one
		for _, dep := range d.edges[id] {
			dfs(dep)
		}

		stack = append(stack, id)
	}

	// Visit all nodes
	for id := range d.nodes {
		dfs(id)
	}

	// Reverse the stack to get dependencies first
	for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		stack[i], stack[j] = stack[j], stack[i]
	}

	return stack
}

// HasCycle checks if the DAG has a cycle
func (d *DefaultDAG) HasCycle() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	visited := make(map[string]int) // 0: unvisited, 1: visiting, 2: visited
	var hasCycle bool

	var dfs func(id string)
	dfs = func(id string) {
		visited[id] = 1 // Mark as visiting

		for _, neighbor := range d.edges[id] {
			if visited[neighbor] == 1 {
				hasCycle = true
				return
			}
			if visited[neighbor] == 0 {
				dfs(neighbor)
			}
		}

		visited[id] = 2 // Mark as visited
	}

	for id := range d.nodes {
		if visited[id] == 0 {
			dfs(id)
		}
	}

	return hasCycle
}

// GetDependencies returns all nodes that the given node depends on
func (d *DefaultDAG) GetDependencies(id string) ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if _, exists := d.nodes[id]; !exists {
		return nil, fmt.Errorf("node %s not found", id)
	}

	deps := []string{}
	visited := make(map[string]bool)

	var collect func(nodeID string)
	collect = func(nodeID string) {
		if visited[nodeID] {
			return
		}
		visited[nodeID] = true

		// Collect incoming edges (dependencies)
		for from, tos := range d.edges {
			for _, to := range tos {
				if to == nodeID && !visited[from] {
					deps = append(deps, from)
					collect(from)
				}
			}
		}
	}

	collect(id)
	return deps, nil
}
