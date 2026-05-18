package discovery

import (
	"testing"
)

func TestNewServiceDiscovery(t *testing.T) {
	d := NewServiceDiscovery()
	if d == nil {
		t.Fatal("discovery should not be nil")
	}
}

func TestRegisterAndGetInstance(t *testing.T) {
	d := NewServiceDiscovery()

	d.RegisterStatic("auth", []*ServiceInstance{
		{ID: "auth-1", Name: "auth", Address: "10.0.1.1", Port: 8080, Weight: 1},
		{ID: "auth-2", Name: "auth", Address: "10.0.1.2", Port: 8080, Weight: 2},
	})

	inst := d.GetInstances("auth")
	if len(inst) != 2 {
		t.Fatalf("expected 2 instances, got %d", len(inst))
	}

	instance, err := d.GetInstance("auth")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if instance == nil {
		t.Fatal("instance should not be nil")
	}
}

func TestGetInstance_NoHealthy(t *testing.T) {
	d := NewServiceDiscovery()
	_, err := d.GetInstance("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent service")
	}
}

func TestMarkUnhealthy(t *testing.T) {
	d := NewServiceDiscovery()

	d.RegisterStatic("catalog", []*ServiceInstance{
		{ID: "cat-1", Name: "catalog", Address: "10.0.2.1", Port: 8080},
	})

	d.MarkUnhealthy("catalog", "cat-1")
	instances := d.GetInstances("catalog")
	if len(instances) != 0 {
		t.Error("expected no healthy instances after marking unhealthy")
	}
}

func TestMarkHealthy(t *testing.T) {
	d := NewServiceDiscovery()

	d.RegisterStatic("catalog", []*ServiceInstance{
		{ID: "cat-1", Name: "catalog", Address: "10.0.2.1", Port: 8080},
	})

	d.MarkUnhealthy("catalog", "cat-1")
	d.MarkHealthy("catalog", "cat-1")
	instances := d.GetInstances("catalog")
	if len(instances) != 1 {
		t.Errorf("expected 1 healthy instance, got %d", len(instances))
	}
}

func TestGetAllServices(t *testing.T) {
	d := NewServiceDiscovery()

	d.RegisterStatic("auth", []*ServiceInstance{{ID: "auth-1", Address: "localhost", Port: 8080}})
	d.RegisterStatic("catalog", []*ServiceInstance{{ID: "cat-1", Address: "localhost", Port: 8081}})

	services := d.GetAllServices()
	if len(services) != 2 {
		t.Errorf("expected 2 services, got %d", len(services))
	}
}

func TestParseServiceTarget(t *testing.T) {
	tests := []struct {
		raw      string
		name     string
		secure   bool
		wantErr  bool
	}{
		{"http://auth:8080", "auth", false, false},
		{"https://catalog:8443", "catalog", true, false},
		{"grpc://inventory:9090", "inventory", false, false},
		{"auth:8080", "auth", false, false},
	}

	for _, tt := range tests {
		target, err := ParseServiceTarget(tt.raw)
		if tt.wantErr {
			if err == nil {
				t.Errorf("expected error for %s", tt.raw)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error for %s: %v", tt.raw, err)
			continue
		}
		if target.Name != tt.name {
			t.Errorf("expected name %s, got %s", tt.name, target.Name)
		}
	}
}

func TestGetInstanceWeighted(t *testing.T) {
	d := NewServiceDiscovery()

	d.RegisterStatic("cart", []*ServiceInstance{
		{ID: "cart-1", Name: "cart", Address: "10.0.3.1", Port: 8080, Weight: 1},
		{ID: "cart-2", Name: "cart", Address: "10.0.3.2", Port: 8080, Weight: 3},
	})

	for i := 0; i < 10; i++ {
		inst, err := d.GetInstanceWeighted("cart")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if inst == nil {
			t.Fatal("instance should not be nil")
		}
	}
}
