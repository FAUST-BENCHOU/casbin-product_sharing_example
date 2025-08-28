package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

// ProductService product service
type ProductService struct {
	enforcer *casbin.Enforcer
	products map[string]*Product
	mu       sync.RWMutex
}

// NewProductService creates a new product service
func NewProductService() (*ProductService, error) {
	// Create model
	m, err := model.NewModelFromString(`
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = (r.sub == r.obj.Owner) || contains(r.obj.SharedUsers, r.sub)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %v", err)
	}

	// Create enforcer
	enforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %v", err)
	}

	// Add custom function
	enforcer.AddFunction("contains", func(args ...interface{}) (interface{}, error) {
		// First parameter should be string array, second parameter is the string to find
		if len(args) != 2 {
			return false, nil
		}
		users, ok := args[0].([]string)
		if !ok {
			return false, nil
		}
		user, ok := args[1].(string)
		if !ok {
			return false, nil
		}
		for _, u := range users {
			if u == user {
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

// ShareProduct shares a product with another user
func (ps *ProductService) ShareProduct(productID, ownerID, targetUserID string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	product, exists := ps.products[productID]
	if !exists {
		return fmt.Errorf("product %s does not exist", productID)
	}

	if product.Owner != ownerID {
		return fmt.Errorf("user %s is not the owner of product %s", ownerID, productID)
	}

	// Check if target user is already in shared users list
	for _, user := range product.SharedUsers {
		if user == targetUserID {
			return fmt.Errorf("user %s is already a shared user", targetUserID)
		}
	}

	return product.AddSharedUser(targetUserID)
}

// CreateProduct creates a new product
func (ps *ProductService) CreateProduct(id, name, ownerID string) (*Product, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.products[id]; exists {
		return nil, fmt.Errorf("product ID %s already exists", id)
	}

	product := NewProduct(id, name, ownerID)
	ps.products[id] = product

	fmt.Printf("Product %s created, owner: %s\n", name, ownerID)
	return product, nil
}

// GetProduct gets a product by ID
func (ps *ProductService) GetProduct(id string) (*Product, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	product, exists := ps.products[id]
	if !exists {
		return nil, fmt.Errorf("product %s does not exist", id)
	}

	return product, nil
}

// CanAccessProduct checks if a user can access a product
func (ps *ProductService) CanAccessProduct(userID, productID, action string) bool {
	product, err := ps.GetProduct(productID)
	if err != nil {
		return false
	}

	// Create a request for Casbin
	request := []interface{}{userID, product, action}
	allowed, err := ps.enforcer.Enforce(request...)
	if err != nil {
		log.Printf("Error checking permission: %v", err)
		return false
	}

	return allowed
}

// UnshareProduct unshares a product from a user
func (ps *ProductService) UnshareProduct(productID, ownerID, targetUserID string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	product, exists := ps.products[productID]
	if !exists {
		return fmt.Errorf("product %s does not exist", productID)
	}

	if product.Owner != ownerID {
		return fmt.Errorf("user %s is not the owner of product %s", ownerID, productID)
	}

	return product.RemoveSharedUser(targetUserID)
}

// ListUserProducts lists all products owned and shared by a user
func (ps *ProductService) ListUserProducts(userID string) ([]*Product, []*Product) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	var owned []*Product
	var shared []*Product

	for _, product := range ps.products {
		if product.Owner == userID {
			owned = append(owned, product)
		} else if product.IsSharedUser(userID) {
			shared = append(shared, product)
		}
	}

	return owned, shared
}

// Close closes the product service
func (ps *ProductService) Close() {
	// Clean up resources if needed
}
