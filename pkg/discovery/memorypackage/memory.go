package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jay-SP/movieapplication/pkg/discovery"
)

type serviceName string
type instanceID string

// Registry defines an in-memory service registry.
type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

// NewRegistry creates a new in-memory service registry instance.
func NewRegistry() *Registry {
	return &Registry{
		serviceAddrs: make(map[serviceName]map[instanceID]*serviceInstance),
	}
}

// Register creates a service record in the registry.
func (r *Registry) Register(ctx context.Context, instanceId instanceID, serviceName serviceName, hostPort string) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok {
		r.serviceAddrs[serviceName] = make(map[instanceID]*serviceInstance)
	}

	r.serviceAddrs[serviceName][instanceId] = &serviceInstance{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}
}

// Deregister removes a service record from the registry.
func (r *Registry) Deregister(ctx context.Context, instanceId instanceID, serviceName serviceName) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return nil // Service not found, no error
	}

	delete(r.serviceAddrs[serviceName], instanceId)
	return nil
}

// ReportHealthyState is a mechanism for reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(instanceID instanceID, serviceName serviceName) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return errors.New("service is not registered yet")
	}

	if _, ok := r.serviceAddrs[serviceName][instanceID]; !ok {
		return errors.New("service instance is not registered yet")
	}

	r.serviceAddrs[serviceName][instanceID].lastActive = time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (r *Registry) ServiceAddresses(ctx context.Context, serviceName serviceName) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok || len(r.serviceAddrs[serviceName]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string
	now := time.Now()

	for _, i := range r.serviceAddrs[serviceName] {
		if i.lastActive.Before(now.Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}
