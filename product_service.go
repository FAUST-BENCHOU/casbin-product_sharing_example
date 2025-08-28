package main

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// ProductService 产品服务
type ProductService struct {
	enforcer *casbin.Enforcer
	products map[string]*Product
	mu       sync.RWMutex
}

// NewProductService 创建产品服务
func NewProductService() (*ProductService, error) {
	// 创建模型
	m, err := model.NewModelFromFile("model.conf")
	if err != nil {
		return nil, fmt.Errorf("创建模型失败: %v", err)
	}

	// 创建enforcer
	enforcer, err := casbin.NewEnforcer(m, fileadapter.NewAdapter(""))
	if err != nil {
		return nil, fmt.Errorf("创建enforcer失败: %v", err)
	}

	// 添加自定义函数
	enforcer.AddFunction("contains", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return false, nil
		}

		// 第一个参数应该是字符串数组，第二个参数是要查找的字符串
		sharedUsers, ok := args[0].([]string)
		if !ok {
			return false, nil
		}

		targetUser, ok := args[1].(string)
		if !ok {
			return false, nil
		}

		// 检查目标用户是否在共享用户列表中
		for _, user := range sharedUsers {
			if user == targetUser {
				return true, nil
			}
		}

		return false, nil
	})

	return &ProductService{
		enforcer: enforcer,
		products: make(map[string]*Product),
	}, nil
}

// CreateProduct 创建产品
func (ps *ProductService) CreateProduct(id, name, ownerID string) (*Product, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.products[id]; exists {
		return nil, fmt.Errorf("产品ID %s 已存在", id)
	}

	product := NewProduct(id, name, ownerID)
	ps.products[id] = product

	fmt.Printf("产品 %s 已创建，所有者: %s\n", name, ownerID)
	return product, nil
}

// GetProduct 获取产品
func (ps *ProductService) GetProduct(id string) (*Product, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	product, exists := ps.products[id]
	if !exists {
		return nil, fmt.Errorf("产品 %s 不存在", id)
	}

	return product, nil
}

// ShareProduct 共享产品给其他用户
func (ps *ProductService) ShareProduct(productID, ownerID, targetUserID string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	product, exists := ps.products[productID]
	if !exists {
		return fmt.Errorf("产品 %s 不存在", productID)
	}

	if product.Owner != ownerID {
		return fmt.Errorf("只有产品所有者可以共享产品")
	}

	if ownerID == targetUserID {
		return fmt.Errorf("不能共享给自己")
	}

	return product.AddSharedUser(targetUserID)
}

// UnshareProduct 取消产品共享
func (ps *ProductService) UnshareProduct(productID, ownerID, targetUserID string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	product, exists := ps.products[productID]
	if !exists {
		return fmt.Errorf("产品 %s 不存在", productID)
	}

	if product.Owner != ownerID {
		return fmt.Errorf("只有产品所有者可以取消共享")
	}

	return product.RemoveSharedUser(targetUserID)
}

// CanAccessProduct 检查用户是否可以访问产品
func (ps *ProductService) CanAccessProduct(userID, productID, action string) bool {
	ps.mu.RLock()
	product, exists := ps.products[productID]
	ps.mu.RUnlock()

	if !exists {
		return false
	}

	// 使用Casbin进行权限检查
	allowed, err := ps.enforcer.Enforce(userID, product, action)
	if err != nil {
		fmt.Printf("权限检查失败: %v\n", err)
		return false
	}

	return allowed
}

// ListUserProducts 列出用户的产品（拥有的和共享的）
func (ps *ProductService) ListUserProducts(userID string) (owned []*Product, shared []*Product) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, product := range ps.products {
		if product.Owner == userID {
			owned = append(owned, product)
		} else if product.IsSharedUser(userID) {
			shared = append(shared, product)
		}
	}

	return owned, shared
}

// Close 关闭服务
func (ps *ProductService) Close() {
	// Casbin enforcer 不需要显式关闭
}
