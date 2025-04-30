package memorypackage

import (
	"context"
	"errors"
	"movieexample.com/pkg/discovery"
	"sync"
	"time"
)

type serviceNameType string
type instanceIDType string

type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceNameType]map[instanceIDType]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[serviceNameType]map[instanceIDType]*serviceInstance{}}
}

func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceNameType(serviceName)]; !ok {
		r.serviceAddrs[serviceNameType(serviceName)] = map[instanceIDType]*serviceInstance{}
	}

	r.serviceAddrs[serviceNameType(serviceName)][instanceIDType(instanceID)] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}

	return nil
}

func (r *Registry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceNameType(serviceName)]; !ok {
		return nil
	}

	delete(r.serviceAddrs[serviceNameType(serviceName)], instanceIDType(instanceID))

	return nil
}

func (r *Registry) ReportHealthyState(instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceNameType(serviceName)]; !ok {
		return errors.New("service is not registered yet")
	}

	_, ok := r.serviceAddrs[serviceNameType(serviceName)][instanceIDType(instanceID)]
	if !ok {
		return errors.New("service instance is not registered yet")
	}

	r.serviceAddrs[serviceNameType(serviceName)][instanceIDType(instanceID)].lastActive = time.Now()

	return nil
}

func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	if _, ok := r.serviceAddrs[serviceNameType(serviceName)]; !ok {
		return nil, errors.New("service is not registered yet")
	}

	if len(r.serviceAddrs[serviceNameType(serviceName)]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var result []string
	for _, service := range r.serviceAddrs[serviceNameType(serviceName)] {
		if service.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		result = append(result, service.hostPort)
	}

	return result, nil
}
