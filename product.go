package main

import (
	"fmt"
	"sync"
)

// Product 表示一个产品
type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Owner       string   `json:"owner"`
	SharedUsers []string `json:"shared_users"`
	mu          sync.RWMutex
}

// NewProduct 创建新产品
func NewProduct(id, name, owner string) *Product {
	return &Product{
		ID:          id,
		Name:        name,
		Owner:       owner,
		SharedUsers: []string{},
	}
}

// AddSharedUser 添加共享用户
func (p *Product) AddSharedUser(userID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// 检查用户是否已经是共享用户
	for _, user := range p.SharedUsers {
		if user == userID {
			return fmt.Errorf("用户 %s 已经是共享用户", userID)
		}
	}
	
	p.SharedUsers = append(p.SharedUsers, userID)
	return nil
}

// RemoveSharedUser 移除共享用户
func (p *Product) RemoveSharedUser(userID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	for i, user := range p.SharedUsers {
		if user == userID {
			p.SharedUsers = append(p.SharedUsers[:i], p.SharedUsers[i+1:]...)
			return nil
		}
	}
	
	return fmt.Errorf("用户 %s 不是共享用户", userID)
}

// IsSharedUser 检查用户是否是共享用户
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

// GetSharedUsers 获取所有共享用户
func (p *Product) GetSharedUsers() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	result := make([]string, len(p.SharedUsers))
	copy(result, p.SharedUsers)
	return result
}
