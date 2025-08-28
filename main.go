package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("=== Casbin 产品共享示例 ===")
	fmt.Println()

	// 创建产品服务
	service, err := NewProductService()
	if err != nil {
		log.Fatalf("创建产品服务失败: %v", err)
	}
	defer service.Close()

	// 创建一些用户
	user1 := "user1@example.com"
	user2 := "user2@example.com"
	user3 := "user3@example.com"

	fmt.Println("1. 创建产品...")

	// 用户1创建产品
	product1, err := service.CreateProduct("prod_001", "iPhone 15", user1)
	if err != nil {
		log.Fatalf("创建产品失败: %v", err)
	}

	product2, err := service.CreateProduct("prod_002", "MacBook Pro", user1)
	if err != nil {
		log.Fatalf("创建产品失败: %v", err)
	}

	// 用户2创建产品
	product3, err := service.CreateProduct("prod_003", "iPad Pro", user2)
	if err != nil {
		log.Fatalf("创建产品失败: %v", err)
	}

	fmt.Println()

	fmt.Println("2. 测试权限控制...")

	// 测试用户1访问自己的产品
	fmt.Printf("用户 %s 访问产品 %s (read): %t\n",
		user1, product1.Name, service.CanAccessProduct(user1, product1.ID, "read"))

	fmt.Printf("用户 %s 访问产品 %s (write): %t\n",
		user1, product1.Name, service.CanAccessProduct(user1, product1.ID, "write"))

	// 测试用户2访问用户1的产品（应该被拒绝）
	fmt.Printf("用户 %s 访问产品 %s (read): %t\n",
		user2, product1.Name, service.CanAccessProduct(user2, product1.ID, "read"))

	// 测试用户3访问用户2的产品（应该被拒绝）
	fmt.Printf("用户 %s 访问产品 %s (read): %t\n",
		user3, product3.Name, service.CanAccessProduct(user3, product3.ID, "read"))

	fmt.Println()

	fmt.Println("3. 共享产品...")

	// 用户1将产品1共享给用户2
	err = service.ShareProduct(product1.ID, user1, user2)
	if err != nil {
		log.Printf("共享产品失败: %v", err)
	} else {
		fmt.Printf("产品 %s 已共享给用户 %s\n", product1.Name, user2)
	}

	// 用户1将产品2共享给用户3
	err = service.ShareProduct(product2.ID, user1, user3)
	if err != nil {
		log.Printf("共享产品失败: %v", err)
	} else {
		fmt.Printf("产品 %s 已共享给用户 %s\n", product2.Name, user3)
	}

	fmt.Println()

	fmt.Println("4. 测试共享后的权限...")

	// 测试用户2访问共享的产品
	fmt.Printf("用户 %s 访问共享产品 %s (read): %t\n",
		user2, product1.Name, service.CanAccessProduct(user2, product1.ID, "read"))

	fmt.Printf("用户 %s 访问共享产品 %s (write): %t\n",
		user2, product1.Name, service.CanAccessProduct(user2, product1.ID, "write"))

	// 测试用户3访问共享的产品
	fmt.Printf("用户 %s 访问共享产品 %s (read): %t\n",
		user3, product2.Name, service.CanAccessProduct(user3, product2.ID, "read"))

	fmt.Println()

	fmt.Println("5. 列出用户的产品...")

	// 列出用户1的产品
	owned1, shared1 := service.ListUserProducts(user1)
	fmt.Printf("用户 %s 拥有的产品: %d 个\n", user1, len(owned1))
	for _, p := range owned1 {
		fmt.Printf("  - %s (ID: %s)\n", p.Name, p.ID)
	}
	fmt.Printf("用户 %s 共享的产品: %d 个\n", user1, len(shared1))

	// 列出用户2的产品
	owned2, shared2 := service.ListUserProducts(user2)
	fmt.Printf("用户 %s 拥有的产品: %d 个\n", user2, len(owned2))
	for _, p := range owned2 {
		fmt.Printf("  - %s (ID: %s)\n", p.Name, p.ID)
	}
	fmt.Printf("用户 %s 共享的产品: %d 个\n", user2, len(shared2))
	for _, p := range shared2 {
		fmt.Printf("  - %s (ID: %s, 所有者: %s)\n", p.Name, p.ID, p.Owner)
	}

	fmt.Println()

	fmt.Println("6. 取消共享...")

	// 用户1取消与用户2的共享
	err = service.UnshareProduct(product1.ID, user1, user2)
	if err != nil {
		log.Printf("取消共享失败: %v", err)
	} else {
		fmt.Printf("产品 %s 已取消与用户 %s 的共享\n", product1.Name, user2)
	}

	// 测试取消共享后的权限
	fmt.Printf("用户 %s 访问已取消共享的产品 %s (read): %t\n",
		user2, product1.Name, service.CanAccessProduct(user2, product1.ID, "read"))

	fmt.Println()
	fmt.Println("=== 示例完成 ===")
}
