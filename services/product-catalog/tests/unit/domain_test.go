package domain

import (
	"testing"
	"time"
	"github.com/shopee-clone/shopee/services/product-catalog/internal/domain"
)

func TestNewProduct(t *testing.T) {
	p := domain.NewProduct("SHOP-001", "Test Product", "Description", "CAT-001", "SGD")
	if p.Name != "Test Product" { t.Errorf("expected Test Product, got %s", p.Name) }
	if p.Status != domain.ProductStatusDraft { t.Errorf("expected draft, got %s", p.Status) }
	if p.ShopID != "SHOP-001" { t.Errorf("expected SHOP-001, got %s", p.ShopID) }
}

func TestProduct_Activate(t *testing.T) {
	p := domain.NewProduct("SHOP-001", "Test", "Desc", "CAT-001", "SGD")
	if err := p.Activate(); err != nil { t.Errorf("unexpected error: %v", err) }
	if p.Status != domain.ProductStatusActive { t.Errorf("expected active, got %s", p.Status) }
}

func TestProduct_Activate_InvalidState(t *testing.T) {
	p := domain.NewProduct("SHOP-001", "Test", "Desc", "CAT-001", "SGD")
	p.Activate()
	if err := p.Activate(); err == nil { t.Error("expected error for double activate") }
}

func TestProduct_Archive(t *testing.T) {
	p := domain.NewProduct("SHOP-001", "Test", "Desc", "CAT-001", "SGD")
	if err := p.Archive(); err != nil { t.Errorf("unexpected error: %v", err) }
	if p.Status != domain.ProductStatusArchived { t.Errorf("expected archived, got %s", p.Status) }
}

func TestProduct_Update(t *testing.T) {
	p := domain.NewProduct("SHOP-001", "Test", "Desc", "CAT-001", "SGD")
	p.Update("Updated Name", "Updated Desc")
	if p.Name != "Updated Name" { t.Errorf("expected Updated Name, got %s", p.Name) }
	if p.Version != 2 { t.Errorf("expected version 2, got %d", p.Version) }
}

func TestNewSKU(t *testing.T) {
	sku := domain.NewSKU("PROD-001", "SKU-001", "color:red", 5000)
	if sku.SKUCode != "SKU-001" { t.Errorf("expected SKU-001, got %s", sku.SKUCode) }
	if sku.Price != 5000 { t.Errorf("expected 5000, got %d", sku.Price) }
	if sku.Status != domain.SKUStatusActive { t.Errorf("expected active, got %s", sku.Status) }
}

func TestNewCategory(t *testing.T) {
	c := domain.NewCategory("", "Electronics", "electronics", 0, 1)
	if c.Name != "Electronics" { t.Errorf("expected Electronics, got %s", c.Name) }
	if c.Level != 0 { t.Errorf("expected level 0, got %d", c.Level) }
	if !c.IsActive { t.Error("expected active") }
}

func TestNewCategory_WithParent(t *testing.T) {
	c := domain.NewCategory("CAT-001", "Phones", "phones", 1, 1)
	if c.ParentID != "CAT-001" { t.Errorf("expected CAT-001, got %s", c.ParentID) }
	if c.Level != 1 { t.Errorf("expected level 1, got %d", c.Level) }
}

func TestNewAttribute(t *testing.T) {
	a := &domain.Attribute{ID: "ATTR-001", CategoryID: "CAT-001", Name: "Color", DisplayName: "Color", AttrType: domain.AttributeTypeSelect, Required: true}
	if a.AttrType != domain.AttributeTypeSelect { t.Errorf("expected select type, got %s", a.AttrType) }
}

func TestNewProductMedia(t *testing.T) {
	m := &domain.ProductMedia{ID: "MEDIA-001", ProductID: "PROD-001", MediaType: domain.MediaTypeImage, URL: "http://img.url", IsPrimary: true}
	if m.MediaType != domain.MediaTypeImage { t.Errorf("expected image type, got %s", m.MediaType) }
	if !m.IsPrimary { t.Error("expected primary") }
}
