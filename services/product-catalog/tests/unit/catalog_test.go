package unit

import (
	"testing"

	"github.com/shopee-clone/shopee/services/product-catalog/internal/domain"
)

func TestNewProduct(t *testing.T) {
	p := domain.NewProduct("shop-1", "Test Product", "A great product", "cat-1", "Nike", "new")
	if p.ShopID != "shop-1" { t.Errorf("expected shop-1, got %s", p.ShopID) }
	if p.Name != "Test Product" { t.Errorf("expected Test Product, got %s", p.Name) }
	if p.Status != domain.ProductStatusDraft { t.Errorf("expected draft, got %s", p.Status) }
	if p.Version != 1 { t.Errorf("expected version 1, got %d", p.Version) }
}

func TestProduct_CanTransitionTo(t *testing.T) {
	tests := []struct{ from, to domain.ProductStatus; expected bool }{
		{domain.ProductStatusDraft, domain.ProductStatusPending, true},
		{domain.ProductStatusDraft, domain.ProductStatusActive, false},
		{domain.ProductStatusPending, domain.ProductStatusActive, true},
		{domain.ProductStatusPending, domain.ProductStatusRejected, true},
		{domain.ProductStatusActive, domain.ProductStatusInactive, true},
		{domain.ProductStatusActive, domain.ProductStatusArchived, true},
		{domain.ProductStatusArchived, domain.ProductStatusActive, false},
	}
	for _, tt := range tests {
		p := &domain.Product{Status: tt.from}
		if p.CanTransitionTo(tt.to) != tt.expected {
			t.Errorf("CanTransitionTo(%s,%s)=%v want %v", tt.from, tt.to, p.CanTransitionTo(tt.to), tt.expected)
		}
	}
}

func TestProduct_IsEditable(t *testing.T) {
	tests := []struct{ status domain.ProductStatus; expected bool }{
		{domain.ProductStatusDraft, true},
		{domain.ProductStatusInactive, true},
		{domain.ProductStatusRejected, true},
		{domain.ProductStatusActive, false},
		{domain.ProductStatusArchived, false},
	}
	for _, tt := range tests {
		p := &domain.Product{Status: tt.status}
		if p.IsEditable() != tt.expected {
			t.Errorf("IsEditable(%s)=%v want %v", tt.status, p.IsEditable(), tt.expected)
		}
	}
}

func TestNewSKU(t *testing.T) {
	sku := domain.NewSKU("prod-1", "SKU-001", "Red / Large", "SGD", 2999)
	if sku.ProductID != "prod-1" { t.Errorf("expected prod-1, got %s", sku.ProductID) }
	if sku.Price != 2999 { t.Errorf("expected 2999, got %d", sku.Price) }
	if !sku.IsAvailable() { t.Error("expected SKU to be available") }
}

func TestSKU_Reserve(t *testing.T) {
	sku := domain.NewSKU("prod-1", "SKU-001", "Red", "SGD", 100)
	sku.Stock = 10
	if err := sku.Reserve(3); err != nil { t.Fatalf("unexpected error: %v", err) }
	if sku.ReservedStock != 3 { t.Errorf("expected reserved 3, got %d", sku.ReservedStock) }
	if sku.AvailableStock() != 7 { t.Errorf("expected available 7, got %d", sku.AvailableStock()) }
}

func TestSKU_ReserveInsufficient(t *testing.T) {
	sku := domain.NewSKU("prod-1", "SKU-001", "Red", "SGD", 100)
	sku.Stock = 5
	if err := sku.Reserve(10); err == nil {
		t.Error("expected error for insufficient stock")
	}
}

func TestNewCategory(t *testing.T) {
	c := domain.NewCategory("Electronics", "electronics", "All electronics", nil, 0)
	if c.Name != "Electronics" { t.Errorf("expected Electronics, got %s", c.Name) }
	if c.Slug != "electronics" { t.Errorf("expected electronics, got %s", c.Slug) }
	if !c.IsRoot() { t.Error("expected root category") }
	if c.Depth != 0 { t.Errorf("expected depth 0, got %d", c.Depth) }
}

func TestCategory_WithParent(t *testing.T) {
	parentID := "parent-1"
	c := domain.NewCategory("Phones", "phones", "Smartphones", &parentID, 1)
	if c.IsRoot() { t.Error("expected non-root category") }
	if c.Depth != 1 { t.Errorf("expected depth 1, got %d", c.Depth) }
}

func TestNewMedia(t *testing.T) {
	m := domain.NewMedia("prod-1", "image", "https://cdn.example.com/img.jpg", 0)
	if m.ProductID != "prod-1" { t.Errorf("expected prod-1, got %s", m.ProductID) }
	if m.Type != domain.MediaTypeImage { t.Errorf("expected image, got %s", m.Type) }
	if m.Status != domain.MediaStatusPending { t.Errorf("expected pending, got %s", m.Status) }
}
