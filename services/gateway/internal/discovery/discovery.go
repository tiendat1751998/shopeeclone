package discovery

import (
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type ServiceInstance struct {
	ID      string
	Name    string
	Address string
	Port    int
	Tags    []string
	Healthy bool
	Weight  int
}

type ServiceDiscovery struct {
	mu        sync.RWMutex
	instances map[string][]*ServiceInstance
	rrIndexes map[string]*uint64
	static    bool
}

func NewServiceDiscovery() *ServiceDiscovery {
	return &ServiceDiscovery{
		instances: make(map[string][]*ServiceInstance),
		rrIndexes: make(map[string]*uint64),
		static:   true,
	}
}

func (d *ServiceDiscovery) RegisterStatic(name string, instances []*ServiceInstance) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i := range instances {
		instances[i].Healthy = true
		if instances[i].Weight <= 0 {
			instances[i].Weight = 1
		}
	}

	d.instances[name] = instances
	idx := uint64(0)
	d.rrIndexes[name] = &idx
}

func (d *ServiceDiscovery) GetInstances(name string) []*ServiceInstance {
	d.mu.RLock()
	defer d.mu.RUnlock()

	instances, exists := d.instances[name]
	if !exists {
		return nil
	}

	healthy := make([]*ServiceInstance, 0, len(instances))
	for _, inst := range instances {
		if inst.Healthy {
			healthy = append(healthy, inst)
		}
	}
	return healthy
}

func (d *ServiceDiscovery) GetInstance(name string) (*ServiceInstance, error) {
	instances := d.GetInstances(name)
	if len(instances) == 0 {
		return nil, fmt.Errorf("no healthy instances for service: %s", name)
	}

	idx := atomic.AddUint64(d.rrIndexes[name], 1) - 1
	return instances[idx%uint64(len(instances))], nil
}

func (d *ServiceDiscovery) GetInstanceWeighted(name string) (*ServiceInstance, error) {
	instances := d.GetInstances(name)
	if len(instances) == 0 {
		return nil, fmt.Errorf("no healthy instances for service: %s", name)
	}

	totalWeight := 0
	for _, inst := range instances {
		totalWeight += inst.Weight
	}

	r := rand.Intn(totalWeight)
	for _, inst := range instances {
		r -= inst.Weight
		if r < 0 {
			return inst, nil
		}
	}

	return instances[0], nil
}

func (d *ServiceDiscovery) MarkUnhealthy(name, instanceID string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	instances, exists := d.instances[name]
	if !exists {
		return
	}

	for _, inst := range instances {
		if inst.ID == instanceID {
			inst.Healthy = false
			return
		}
	}
}

func (d *ServiceDiscovery) MarkHealthy(name, instanceID string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	instances, exists := d.instances[name]
	if !exists {
		return
	}

	for _, inst := range instances {
		if inst.ID == instanceID {
			inst.Healthy = true
			return
		}
	}
}

func (d *ServiceDiscovery) GetAllServices() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	services := make([]string, 0, len(d.instances))
	for name := range d.instances {
		services = append(services, name)
	}
	sort.Strings(services)
	return services
}

type ServiceTarget struct {
	Name    string
	Address string
	Secure  bool
	Timeout time.Duration
}

func ParseServiceTarget(raw string) (*ServiceTarget, error) {
	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid service target: %w", err)
	}

	target := &ServiceTarget{
		Name:    u.Hostname(),
		Address: u.Host,
		Secure:  u.Scheme == "https",
	}

	if u.Scheme == "grpc" || u.Scheme == "grpcs" {
		target.Secure = u.Scheme == "grpcs"
	}

	return target, nil
}
