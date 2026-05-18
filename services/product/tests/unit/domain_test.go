package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestProduct_IsAvailable(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		expected bool
	}{
		{
			name:     "active product is available",
			product:  Product{Status: ProductStatusActive},
			expected: true,
		},
		{
			name:     "draft product is not available",
			product:  Product{Status: ProductStatusDraft},
			expected: false,
		},
		{
			name:     "deleted product is not available",
			product:  Product{Status: ProductStatusActive, DeletedAt: &time.Time{}},
			expected: false,
		},
		{
			name:     "inactive product is not available",
			product:  Product{Status: ProductStatusInactive},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.product.IsAvailable(); got != tt.expected {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProduct_IsListable(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		expected bool
	}{
		{"active is listable", Product{Status: ProductStatusActive}, true},
		{"draft is not listable", Product{Status: ProductStatusDraft}, false},
		{"inactive is not listable", Product{Status: ProductStatusInactive}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.product.IsListable(); got != tt.expected {
				t.Errorf("IsListable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProduct_HasStock(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		expected bool
	}{
		{
			name: "has stock",
			product: Product{SKUs: []SKU{
				{Status: SKUStatusActive, Stock: 10},
			}},
			expected: true,
		},
		{
			name: "no stock",
			product: Product{SKUs: []SKU{
				{Status: SKUStatusActive, Stock: 0},
			}},
			expected: false,
		},
		{
			name: "inactive sku with stock",
			product: Product{SKUs: []SKU{
				{Status: SKUStatusInactive, Stock: 10},
			}},
			expected: false,
		},
		{
			name:     "no skus",
			product:  Product{SKUs: []SKU{}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.product.HasStock(); got != tt.expected {
				t.Errorf("HasStock() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProduct_PrimaryImage(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		product  Product
		expected string
	}{
		{
			name: "has primary image",
			product: Product{Images: []ProductImage{
				{URL: "secondary.jpg", IsPrimary: false},
				{URL: "primary.jpg", IsPrimary: true},
			}},
			expected: "primary.jpg",
		},
		{
			name: "no primary, returns first",
			product: Product{Images: []ProductImage{
				{URL: "first.jpg", IsPrimary: false},
			}},
			expected: "first.jpg",
		},
		{
			name:     "no images",
			product:  Product{Images: []ProductImage{}},
			expected: "",
		},
		{
			name:     "nil images",
			product:  Product{},
			expected: "",
		},
	}

	_ = now
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.product.PrimaryImage(); got != tt.expected {
				t.Errorf("PrimaryImage() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSKU_EffectivePrice(t *testing.T) {
	tests := []struct {
		name     string
		sku      SKU
		expected float64
	}{
		{"sale price set", SKU{Price: 100, SalePrice: 80}, 80},
		{"no sale price", SKU{Price: 100, SalePrice: 0}, 100},
		{"sale price higher", SKU{Price: 100, SalePrice: 120}, 120},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sku.EffectivePrice(); got != tt.expected {
				t.Errorf("EffectivePrice() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSKU_IsAvailable(t *testing.T) {
	tests := []struct {
		name     string
		sku      SKU
		expected bool
	}{
		{"active with stock", SKU{Status: SKUStatusActive, Stock: 5}, true},
		{"active no stock", SKU{Status: SKUStatusActive, Stock: 0}, false},
		{"inactive with stock", SKU{Status: SKUStatusInactive, Stock: 5}, false},
		{"out of stock", SKU{Status: SKUStatusOutOfStock, Stock: 0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sku.IsAvailable(); got != tt.expected {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSKU_Volume(t *testing.T) {
	sku := SKU{Length: 10, Width: 5, Height: 2}
	expected := float64(100)
	if got := sku.Volume(); got != expected {
		t.Errorf("Volume() = %v, want %v", got, expected)
	}
}

func TestProductFilter_Normalize(t *testing.T) {
	tests := []struct {
		name     string
		input    ProductFilter
		expected ProductFilter
	}{
		{
			name:     "defaults",
			input:    ProductFilter{},
			expected: ProductFilter{Page: 1, Size: 20, SortBy: SortByCreatedAt, SortOrder: SortDesc},
		},
		{
			name:     "page zero becomes 1",
			input:    ProductFilter{Page: 0},
			expected: ProductFilter{Page: 1, Size: 20, SortBy: SortByCreatedAt, SortOrder: SortDesc},
		},
		{
			name:     "size capped at 100",
			input:    ProductFilter{Size: 200},
			expected: ProductFilter{Page: 1, Size: 100, SortBy: SortByCreatedAt, SortOrder: SortDesc},
		},
		{
			name:     "custom values preserved",
			input:    ProductFilter{Page: 3, Size: 50, SortBy: SortByPrice, SortOrder: SortAsc},
			expected: ProductFilter{Page: 3, Size: 50, SortBy: SortByPrice, SortOrder: SortAsc},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.Normalize()
			if tt.input.Page != tt.expected.Page || tt.input.Size != tt.expected.Size ||
				tt.input.SortBy != tt.expected.SortBy || tt.input.SortOrder != tt.expected.SortOrder {
				t.Errorf("Normalize() = %+v, want %+v", tt.input, tt.expected)
			}
		})
	}
}

func TestProductList_TotalPages(t *testing.T) {
	tests := []struct {
		name     string
		list     ProductList
		expected int
	}{
		{"exact pages", ProductList{Total: 100, Size: 20}, 5},
		{"partial page", ProductList{Total: 105, Size: 20}, 6},
		{"empty", ProductList{Total: 0, Size: 20}, 0},
		{"zero size", ProductList{Total: 100, Size: 0}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.TotalPages(); got != tt.expected {
				t.Errorf("TotalPages() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProductList_HasNext(t *testing.T) {
	list := ProductList{Total: 100, Page: 1, Size: 20}
	if !list.HasNext() {
		t.Error("HasNext() should be true for page 1 of 5")
	}
	list.Page = 5
	if list.HasNext() {
		t.Error("HasNext() should be false for last page")
	}
}

func TestProductList_HasPrevious(t *testing.T) {
	list := ProductList{Total: 100, Page: 1, Size: 20}
	if list.HasPrevious() {
		t.Error("HasPrevious() should be false for page 1")
	}
	list.Page = 3
	if !list.HasPrevious() {
		t.Error("HasPrevious() should be true for page 3")
	}
}

func TestCategory_IsRoot(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		expected bool
	}{
		{"no parent", Category{ParentID: ""}, true},
		{"has parent", Category{ParentID: "CAT-001"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.IsRoot(); got != tt.expected {
				t.Errorf("IsRoot() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCategoryTreeNode_WalkDepthFirst(t *testing.T) {
	root := &CategoryTreeNode{
		Category: Category{CategoryID: "root", Name: "Root"},
		Children: []*CategoryTreeNode{
			{Category: Category{CategoryID: "child1", Name: "Child 1"}},
			{Category: Category{CategoryID: "child2", Name: "Child 2"}},
		},
	}

	var visited []string
	root.WalkDepthFirst(func(node *CategoryTreeNode) {
		visited = append(visited, node.Category.CategoryID)
	})

	expected := []string{"root", "child1", "child2"}
	if len(visited) != len(expected) {
		t.Fatalf("WalkDepthFirst visited %d nodes, expected %d", len(visited), len(expected))
	}
	for i, id := range expected {
		if visited[i] != id {
			t.Errorf("WalkDepthFirst[%d] = %s, want %s", i, visited[i], id)
		}
	}
}

func TestCategoryTreeNode_AllCategoryIDs(t *testing.T) {
	root := &CategoryTreeNode{
		Category: Category{CategoryID: "root"},
		Children: []*CategoryTreeNode{
			{Category: Category{CategoryID: "child1"}},
			{Category: Category{CategoryID: "child2"}},
		},
	}

	ids := root.AllCategoryIDs()
	if len(ids) != 3 {
		t.Errorf("AllCategoryIDs() returned %d IDs, expected 3", len(ids))
	}
}

func TestCategoryTreeNode_FindByCategoryID(t *testing.T) {
	root := &CategoryTreeNode{
		Category: Category{CategoryID: "root"},
		Children: []*CategoryTreeNode{
			{Category: Category{CategoryID: "child1"}},
			{Category: Category{CategoryID: "child2"}},
		},
	}

	found := root.FindByCategoryID("child1")
	if found == nil || found.Category.CategoryID != "child1" {
		t.Error("FindByCategoryID(child1) should find the node")
	}

	notFound := root.FindByCategoryID("nonexistent")
	if notFound != nil {
		t.Error("FindByCategoryID(nonexistent) should return nil")
	}
}

func TestDomainErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrProductNotFound", ErrProductNotFound},
		{"ErrSKUNotFound", ErrSKUNotFound},
		{"ErrCategoryNotFound", ErrCategoryNotFound},
		{"ErrAttributeNotFound", ErrAttributeNotFound},
		{"ErrDuplicateProduct", ErrDuplicateProduct},
		{"ErrDuplicateSKU", ErrDuplicateSKU},
		{"ErrProductNotActive", ErrProductNotActive},
		{"ErrInvalidPrice", ErrInvalidPrice},
		{"ErrInvalidStock", ErrInvalidStock},
		{"ErrInvalidCategory", ErrInvalidCategory},
		{"ErrUnauthorizedOperation", ErrUnauthorizedOperation},
		{"ErrProductLocked", ErrProductLocked},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s should not be nil", tt.name)
			}
			if tt.err.Error() == "" {
				t.Errorf("%s.Error() should not be empty", tt.name)
			}
		})
	}
}

func TestIsDomainError(t *testing.T) {
	if !IsDomainError(ErrProductNotFound, "PRODUCT_NOT_FOUND") {
		t.Error("IsDomainError should return true for matching domain error")
	}
	if IsDomainError(ErrProductNotFound, "WRONG_CODE") {
		t.Error("IsDomainError should return false for non-matching code")
	}
	if IsDomainError(nil, "PRODUCT_NOT_FOUND") {
		t.Error("IsDomainError should return false for nil")
	}
}

func TestIsNotFound(t *testing.T) {
	if !IsNotFound(ErrProductNotFound) {
		t.Error("IsNotFound should return true for ErrProductNotFound")
	}
	if !IsNotFound(ErrCategoryNotFound) {
		t.Error("IsNotFound should return true for ErrCategoryNotFound")
	}
	if IsNotFound(ErrInvalidPrice) {
		t.Error("IsNotFound should return false for non-not-found errors")
	}
}

func TestNewProductCreatedEvent(t *testing.T) {
	product := &Product{SPUID: "SPU-001", Title: "Test Product"}
	event := NewProductCreatedEvent(product)

	if event.Type != EventTypeProductCreated {
		t.Errorf("EventType = %s, want %s", event.Type, EventTypeProductCreated)
	}
	if event.Product != product {
		t.Error("Product should match")
	}
	if event.Timestamp.IsZero() {
		t.Error("Timestamp should be set")
	}
}

func TestNewProductUpdatedEvent(t *testing.T) {
	product := &Product{SPUID: "SPU-001", Title: "Updated"}
	changedFields := []string{"title", "description"}
	event := NewProductUpdatedEvent(product, changedFields)

	if event.Type != EventTypeProductUpdated {
		t.Errorf("EventType = %s, want %s", event.Type, EventTypeProductUpdated)
	}
	if len(event.ChangedFields) != 2 {
		t.Errorf("ChangedFields length = %d, want 2", len(event.ChangedFields))
	}
}

func TestNewCategoryUpdatedEvent(t *testing.T) {
	category := &Category{CategoryID: "CAT-001", Name: "Electronics"}
	event := NewCategoryUpdatedEvent(category)

	if event.Type != EventTypeCategoryUpdated {
		t.Errorf("EventType = %s, want %s", event.Type, EventTypeCategoryUpdated)
	}
	if event.Category != category {
		t.Error("Category should match")
	}
}

func TestEventMarshal(t *testing.T) {
	product := &Product{SPUID: "SPU-001", Title: "Test"}
	event := NewProductCreatedEvent(product)

	data, err := event.Marshal()
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	if len(data) == 0 {
		t.Error("Marshal() should return non-empty data")
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Errorf("Marshal() produced invalid JSON: %v", err)
	}
}
