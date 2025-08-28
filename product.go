package main

import (
	"fmt"
	"sync"
)

// Product represents a product
type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Owner       string   `json:"owner"`
	SharedUsers []string `json:"shared_users"`
	mu          sync.RWMutex
}

// NewProduct creates a new product
func NewProduct(id, name, owner string) *Product {
	return &Product{
		ID:          id,
		Name:        name,
		Owner:       owner,
		SharedUsers: []string{},
	}
}

// AddSharedUser adds a shared user
func (p *Product) AddSharedUser(userID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Check if user is already a shared user
	for _, user := range p.SharedUsers {
		if user == userID {
			return fmt.Errorf("user %s is already a shared user", userID)
		}
	}
	
	p.SharedUsers = append(p.SharedUsers, userID)
	return nil
}

// RemoveSharedUser removes a shared user
func (p *Product) RemoveSharedUser(userID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	for i, user := range p.SharedUsers {
		if user == userID {
			p.SharedUsers = append(p.SharedUsers[:i], p.SharedUsers[i+1:]...)
			return nil
		}
	}
	
	return fmt.Errorf("user %s is not a shared user", userID)
}

// IsSharedUser checks if a user is a shared user
func (p *Product) IsSharedUser(userID string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	for _, user := range p.SharedUsers {
		if user == userID {
			return true
		}
	}
	return false
}

// GetSharedUsers gets all shared users
func (p *Product) GetSharedUsers() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	result := make([]string, len(p.SharedUsers))
	copy(result, p.SharedUsers)
	return result
}
