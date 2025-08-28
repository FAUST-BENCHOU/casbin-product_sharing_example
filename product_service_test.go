package main

import (
	"testing"
)

func TestProductCreation(t *testing.T) {
	service, err := NewProductService()
	if err != nil {
		t.Fatalf("创建产品服务失败: %v", err)
	}
	defer service.Close()

	// 测试创建产品
	product, err := service.CreateProduct("test_001", "测试产品", "test_user")
	if err != nil {
		t.Fatalf("创建产品失败: %v", err)
	}

	if product.ID != "test_001" {
		t.Errorf("产品ID不匹配，期望: test_001, 实际: %s", product.ID)
	}

	if product.Name != "测试产品" {
		t.Errorf("产品名称不匹配，期望: 测试产品, 实际: %s", product.Name)
	}

	if product.Owner != "test_user" {
		t.Errorf("产品所有者不匹配，期望: test_user, 实际: %s", product.Owner)
	}
}

func TestProductSharing(t *testing.T) {
	service, err := NewProductService()
	if err != nil {
		t.Fatalf("创建产品服务失败: %v", err)
	}
	defer service.Close()

	// 创建产品
	product, err := service.CreateProduct("test_002", "共享产品", "owner_user")
	if err != nil {
		t.Fatalf("创建产品失败: %v", err)
	}

	// 测试共享产品
	err = service.ShareProduct(product.ID, "owner_user", "shared_user")
	if err != nil {
		t.Fatalf("共享产品失败: %v", err)
	}

	// 验证共享用户列表
	sharedUsers := product.GetSharedUsers()
	if len(sharedUsers) != 1 {
		t.Errorf("共享用户数量不匹配，期望: 1, 实际: %d", len(sharedUsers))
	}

	if sharedUsers[0] != "shared_user" {
		t.Errorf("共享用户不匹配，期望: shared_user, 实际: %s", sharedUsers[0])
	}

	// 测试重复共享
	err = service.ShareProduct(product.ID, "owner_user", "shared_user")
	if err == nil {
		t.Error("重复共享应该失败")
	}
}

func TestProductUnsharing(t *testing.T) {
	service, err := NewProductService()
	if err != nil {
		t.Fatalf("创建产品服务失败: %v", err)
	}
	defer service.Close()

	// 创建产品
	product, err := service.CreateProduct("test_003", "取消共享产品", "owner_user")
	if err != nil {
		t.Fatalf("创建产品失败: %v", err)
	}

	// 先共享产品
	err = service.ShareProduct(product.ID, "owner_user", "shared_user")
	if err != nil {
		t.Fatalf("共享产品失败: %v", err)
	}

	// 测试取消共享
	err = service.UnshareProduct(product.ID, "owner_user", "shared_user")
	if err != nil {
		t.Fatalf("取消共享失败: %v", err)
	}

	// 验证共享用户列表
	sharedUsers := product.GetSharedUsers()
	if len(sharedUsers) != 0 {
		t.Errorf("共享用户数量不匹配，期望: 0, 实际: %d", len(sharedUsers))
	}
}

func TestAccessControl(t *testing.T) {
	service, err := NewProductService()
	if err != nil {
		t.Fatalf("创建产品服务失败: %v", err)
	}
	defer service.Close()

	// 创建产品
	product, err := service.CreateProduct("test_004", "权限测试产品", "owner_user")
	if err != nil {
		t.Fatalf("创建产品失败: %v", err)
	}

	// 测试所有者访问权限
	if !service.CanAccessProduct("owner_user", product.ID, "read") {
		t.Error("所有者应该能够读取产品")
	}

	if !service.CanAccessProduct("owner_user", product.ID, "write") {
		t.Error("所有者应该能够写入产品")
	}

	// 测试非所有者访问权限（应该被拒绝）
	if service.CanAccessProduct("other_user", product.ID, "read") {
		t.Error("非所有者不应该能够读取产品")
	}

	// 共享产品
	err = service.ShareProduct(product.ID, "owner_user", "shared_user")
	if err != nil {
		t.Fatalf("共享产品失败: %v", err)
	}

	// 测试共享用户访问权限
	if !service.CanAccessProduct("shared_user", product.ID, "read") {
		t.Error("共享用户应该能够读取产品")
	}

	if !service.CanAccessProduct("shared_user", product.ID, "write") {
		t.Error("共享用户应该能够写入产品")
	}
}

func TestListUserProducts(t *testing.T) {
	service, err := NewProductService()
	if err != nil {
		t.Fatalf("创建产品服务失败: %v", err)
	}
	defer service.Close()

	// 创建产品
	product1, err := service.CreateProduct("test_005", "产品1", "user1")
	if err != nil {
		t.Fatalf("创建产品失败: %v", err)
	}

	_, err = service.CreateProduct("test_006", "产品2", "user2")
	if err != nil {
		t.Fatalf("创建产品失败: %v", err)
	}

	// 共享产品
	err = service.ShareProduct(product1.ID, "user1", "user2")
	if err != nil {
		t.Fatalf("共享产品失败: %v", err)
	}

	// 测试用户1的产品列表
	owned1, shared1 := service.ListUserProducts("user1")
	if len(owned1) != 1 {
		t.Errorf("用户1拥有的产品数量不匹配，期望: 1, 实际: %d", len(owned1))
	}

	if len(shared1) != 0 {
		t.Errorf("用户1共享的产品数量不匹配，期望: 0, 实际: %d", len(shared1))
	}

	// 测试用户2的产品列表
	owned2, shared2 := service.ListUserProducts("user2")
	if len(owned2) != 1 {
		t.Errorf("用户2拥有的产品数量不匹配，期望: 1, 实际: %d", len(owned2))
	}

	if len(shared2) != 1 {
		t.Errorf("用户2共享的产品数量不匹配，期望: 1, 实际: %d", len(shared2))
	}
}
