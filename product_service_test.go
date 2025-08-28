package main

import (
	"testing"
)

func TestCreateProduct(t *testing.T) {
	service, err := NewProductService()
	if err != nil {
		t.Fatalf("Failed to create product service: %v", err)
	}
	defer service.Close()

	// Test creating product
	product, err := service.CreateProduct("test_001", "Test Product", "test_user")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Verify product properties
	if product.ID != "test_001" {
		t.Errorf("Product ID mismatch, expected: test_001, got: %s", product.ID)
	}

	if product.Name != "Test Product" {
		t.Errorf("Product name mismatch, expected: Test Product, got: %s", product.Name)
	}

	if product.Owner != "test_user" {
		t.Errorf("Product owner mismatch, expected: test_user, got: %s", product.Owner)
	}
}

func TestShareProduct(t *testing.T) {
	service, err := NewProductService()
	if err != nil {
		t.Fatalf("Failed to create product service: %v", err)
	}
	defer service.Close()

	// Create product
	product, err := service.CreateProduct("test_002", "Shared Product", "owner_user")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Test sharing product
	err = service.ShareProduct("test_002", "owner_user", "shared_user")
	if err != nil {
		t.Fatalf("Failed to share product: %v", err)
	}

	// Verify shared users list
	if !product.IsSharedUser("shared_user") {
		t.Error("Shared user not found in shared users list")
	}
}

func TestProductUnsharing(t *testing.T) {
	service, err := NewProductService()
	if err != nil {
		t.Fatalf("Failed to create product service: %v", err)
	}
	defer service.Close()

	// Create product
	product, err := service.CreateProduct("test_003", "Unshare Product", "owner_user")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Share product first
	err = service.ShareProduct(product.ID, "owner_user", "shared_user")
	if err != nil {
		t.Fatalf("Failed to share product: %v", err)
	}

	// Test unsharing product
	err = service.UnshareProduct(product.ID, "owner_user", "shared_user")
	if err != nil {
		t.Fatalf("Failed to unshare product: %v", err)
	}

	// Verify shared users list
	sharedUsers := product.GetSharedUsers()
	if len(sharedUsers) != 0 {
		t.Errorf("Shared users count mismatch, expected: 0, got: %d", len(sharedUsers))
	}
}

func TestAccessControl(t *testing.T) {
	service, err := NewProductService()
	if err != nil {
		t.Fatalf("Failed to create product service: %v", err)
	}
	defer service.Close()

	// Create product
	product, err := service.CreateProduct("test_004", "Permission Test Product", "owner_user")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Test owner access permissions
	if !service.CanAccessProduct("owner_user", product.ID, "read") {
		t.Error("Owner should be able to read product")
	}

	if !service.CanAccessProduct("owner_user", product.ID, "write") {
		t.Error("Owner should be able to write product")
	}

	// Test non-owner access permissions (should be denied)
	if service.CanAccessProduct("other_user", product.ID, "read") {
		t.Error("Non-owner should not be able to read product")
	}

	// Share product
	err = service.ShareProduct(product.ID, "owner_user", "shared_user")
	if err != nil {
		t.Fatalf("Failed to share product: %v", err)
	}

	// Test shared user access permissions
	if !service.CanAccessProduct("shared_user", product.ID, "read") {
		t.Error("Shared user should be able to read product")
	}

	if !service.CanAccessProduct("shared_user", product.ID, "write") {
		t.Error("Shared user should be able to write product")
	}
}

func TestListUserProducts(t *testing.T) {
	service, err := NewProductService()
	if err != nil {
		t.Fatalf("Failed to create product service: %v", err)
	}
	defer service.Close()

	// Create products
	product1, err := service.CreateProduct("test_005", "Product 1", "user1")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	_, err = service.CreateProduct("test_006", "Product 2", "user2")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Share product
	err = service.ShareProduct(product1.ID, "user1", "user2")
	if err != nil {
		t.Fatalf("Failed to share product: %v", err)
	}

	// Test user1's product list
	owned1, shared1 := service.ListUserProducts("user1")
	if len(owned1) != 1 {
		t.Errorf("User1 owned products count mismatch, expected: 1, got: %d", len(owned1))
	}

	if len(shared1) != 0 {
		t.Errorf("User1 shared products count mismatch, expected: 0, got: %d", len(shared1))
	}

	// Test user2's product list
	owned2, shared2 := service.ListUserProducts("user2")
	if len(owned2) != 1 {
		t.Errorf("User2 owned products count mismatch, expected: 1, got: %d", len(owned2))
	}

	if len(shared2) != 1 {
		t.Errorf("User2 shared products count mismatch, expected: 1, got: %d", len(shared2))
	}
}
