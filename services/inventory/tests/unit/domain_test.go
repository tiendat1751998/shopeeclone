package domain

import (
	"testing"
	"time"

	"github.com/shopee-clone/shopee/services/inventory/internal/domain"
)

func TestNewStock(t *testing.T) {
	stock := domain.NewStock("SKU-001", "WH-001", 100)
	if stock.SKU != "SKU-001" {
		t.Errorf("expected SKU-001, got %s", stock.SKU)
	}
	if stock.WarehouseID != "WH-001" {
		t.Errorf("expected WH-001, got %s", stock.WarehouseID)
	}
	if stock.Quantity != 100 {
		t.Errorf("expected quantity 100, got %d", stock.Quantity)
	}
	if stock.Available != 100 {
		t.Errorf("expected available 100, got %d", stock.Available)
	}
	if stock.Reserved != 0 {
		t.Errorf("expected reserved 0, got %d", stock.Reserved)
	}
	if stock.Version != 1 {
		t.Errorf("expected version 1, got %d", stock.Version)
	}
}

func TestStock_Reserve(t *testing.T) {
	tests := []struct {
		name      string
		quantity  int64
		reserve   int64
		wantErr   bool
		wantAvail int64
		wantRes   int64
	}{
		{"valid reserve", 100, 10, false, 90, 10},
		{"reserve all", 100, 100, false, 0, 100},
		{"over reserve", 100, 101, true, 100, 0},
		{"zero reserve", 100, 0, true, 100, 0},
		{"negative reserve", 100, -1, true, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stock := domain.NewStock("SKU-001", "WH-001", tt.quantity)
			err := stock.Reserve(tt.reserve)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reserve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if stock.Available != tt.wantAvail {
					t.Errorf("Available = %d, want %d", stock.Available, tt.wantAvail)
				}
				if stock.Reserved != tt.wantRes {
					t.Errorf("Reserved = %d, want %d", stock.Reserved, tt.wantRes)
				}
			}
		})
	}
}

func TestStock_Release(t *testing.T) {
	stock := domain.NewStock("SKU-001", "WH-001", 100)
	stock.Reserve(50)

	err := stock.Release(20)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stock.Available != 70 {
		t.Errorf("expected available 70, got %d", stock.Available)
	}
	if stock.Reserved != 30 {
		t.Errorf("expected reserved 30, got %d", stock.Reserved)
	}
}

func TestStock_Release_OverRelease(t *testing.T) {
	stock := domain.NewStock("SKU-001", "WH-001", 100)
	stock.Reserve(10)

	err := stock.Release(20)
	if err == nil {
		t.Error("expected error for over-release")
	}
}

func TestStock_Confirm(t *testing.T) {
	stock := domain.NewStock("SKU-001", "WH-001", 100)
	stock.Reserve(30)

	err := stock.Confirm(30)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stock.Quantity != 70 {
		t.Errorf("expected quantity 70, got %d", stock.Quantity)
	}
	if stock.Reserved != 0 {
		t.Errorf("expected reserved 0, got %d", stock.Reserved)
	}
}

func TestStock_Increase(t *testing.T) {
	stock := domain.NewStock("SKU-001", "WH-001", 100)
	err := stock.Increase(50)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stock.Quantity != 150 {
		t.Errorf("expected quantity 150, got %d", stock.Quantity)
	}
	if stock.Available != 150 {
		t.Errorf("expected available 150, got %d", stock.Available)
	}
}

func TestStock_Decrease(t *testing.T) {
	stock := domain.NewStock("SKU-001", "WH-001", 100)
	err := stock.Decrease(30)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stock.Quantity != 70 {
		t.Errorf("expected quantity 70, got %d", stock.Quantity)
	}
	if stock.Available != 70 {
		t.Errorf("expected available 70, got %d", stock.Available)
	}
}

func TestNewReservation(t *testing.T) {
	res := domain.NewReservation("RES-001", "USER-001", "SKU-001", "WH-001", 5, 15*time.Minute)
	if res.Status != domain.ReservationStatusPending {
		t.Errorf("expected pending status, got %s", res.Status)
	}
	if res.Quantity != 5 {
		t.Errorf("expected quantity 5, got %d", res.Quantity)
	}
	if res.IsExpired() {
		t.Error("new reservation should not be expired")
	}
}

func TestReservation_Confirm(t *testing.T) {
	res := domain.NewReservation("RES-001", "USER-001", "SKU-001", "WH-001", 5, 15*time.Minute)
	if err := res.Confirm(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if res.Status != domain.ReservationStatusConfirmed {
		t.Errorf("expected confirmed status, got %s", res.Status)
	}
}

func TestReservation_Confirm_InvalidTransition(t *testing.T) {
	res := domain.NewReservation("RES-001", "USER-001", "SKU-001", "WH-001", 5, 15*time.Minute)
	res.Confirm()
	err := res.Confirm()
	if err == nil {
		t.Error("expected error for double confirm")
	}
}

func TestReservation_Release(t *testing.T) {
	res := domain.NewReservation("RES-001", "USER-001", "SKU-001", "WH-001", 5, 15*time.Minute)
	if err := res.Release(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if res.Status != domain.ReservationStatusReleased {
		t.Errorf("expected released status, got %s", res.Status)
	}
}

func TestReservation_Expire(t *testing.T) {
	res := domain.NewReservation("RES-001", "USER-001", "SKU-001", "WH-001", 5, 15*time.Minute)
	if err := res.Expire(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if res.Status != domain.ReservationStatusExpired {
		t.Errorf("expected expired status, got %s", res.Status)
	}
}

func TestFlashSaleInventory_IsAvailable(t *testing.T) {
	now := time.Now()
	fs := &domain.FlashSaleInventory{
		FlashSaleID:   "FS-001",
		SKU:           "SKU-001",
		TotalStock:    100,
		ReservedStock: 30,
		SoldStock:     20,
		StartTime:     now.Add(-1 * time.Hour),
		EndTime:       now.Add(1 * time.Hour),
		IsActive:      true,
	}

	if !fs.IsAvailable() {
		t.Error("flash sale should be available")
	}
	if fs.AvailableStock() != 50 {
		t.Errorf("expected available stock 50, got %d", fs.AvailableStock())
	}
}

func TestFlashSaleInventory_IsAvailable_Inactive(t *testing.T) {
	now := time.Now()
	fs := &domain.FlashSaleInventory{
		IsActive:  false,
		StartTime: now.Add(-1 * time.Hour),
		EndTime:   now.Add(1 * time.Hour),
	}
	if fs.IsAvailable() {
		t.Error("inactive flash sale should not be available")
	}
}

func TestStock_ReserveFlashSale(t *testing.T) {
	stock := domain.NewStock("SKU-001", "WH-001", 100)
	stock.FlashSaleStock = 50

	err := stock.ReserveFlashSale(10)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stock.FlashSaleReserved != 10 {
		t.Errorf("expected flash_sale_reserved 10, got %d", stock.FlashSaleReserved)
	}
	if stock.Available != 90 {
		t.Errorf("expected available 90, got %d", stock.Available)
	}
}

func TestStock_ReserveFlashSale_Insufficient(t *testing.T) {
	stock := domain.NewStock("SKU-001", "WH-001", 100)
	stock.FlashSaleStock = 5

	err := stock.ReserveFlashSale(10)
	if err == nil {
		t.Error("expected error for insufficient flash sale stock")
	}
}

// Concurrency test: simulate concurrent reservations
func TestStock_ConcurrentReserve(t *testing.T) {
	stock := domain.NewStock("SKU-001", "WH-001", 1000)

	done := make(chan error, 100)
	for i := 0; i < 100; i++ {
		go func() {
			err := stock.Reserve(1)
			done <- err
		}()
	}

	successCount := 0
	for i := 0; i < 100; i++ {
		if err := <-done; err == nil {
			successCount++
		}
	}

	if successCount != 100 {
		t.Errorf("expected 100 successful reserves, got %d", successCount)
	}
	if stock.Available != 900 {
		t.Errorf("expected available 900, got %d", stock.Available)
	}
	if stock.Reserved != 100 {
		t.Errorf("expected reserved 100, got %d", stock.Reserved)
	}
}
