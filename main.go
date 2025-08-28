package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("=== Casbin Product Sharing Example ===")
	fmt.Println()

	// Create product service
	service, err := NewProductService()
	if err != nil {
		log.Fatalf("Failed to create product service: %v", err)
	}
	defer service.Close()

	// Create some users
	user1 := "user1@example.com"
	user2 := "user2@example.com"
	user3 := "user3@example.com"

	fmt.Println("1. Creating products...")

	// User1 creates products
	product1, err := service.CreateProduct("prod_001", "iPhone 15", user1)
	if err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}

	product2, err := service.CreateProduct("prod_002", "MacBook Pro", user1)
	if err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}

	// User2 creates product
	product3, err := service.CreateProduct("prod_003", "iPad Pro", user2)
	if err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}

	fmt.Println()

	fmt.Println("2. Testing permission control...")

	// Test user1 accessing their own product
	fmt.Printf("User %s accessing product %s (read): %t\n",
		user1, product1.Name, service.CanAccessProduct(user1, product1.ID, "read"))

	fmt.Printf("User %s accessing product %s (write): %t\n",
		user1, product1.Name, service.CanAccessProduct(user1, product1.ID, "write"))

	// Test user2 accessing user1's product (should be denied)
	fmt.Printf("User %s accessing product %s (read): %t\n",
		user2, product1.Name, service.CanAccessProduct(user2, product1.ID, "read"))

	// Test user3 accessing user2's product (should be denied)
	fmt.Printf("User %s accessing product %s (read): %t\n",
		user3, product3.Name, service.CanAccessProduct(user3, product3.ID, "read"))

	fmt.Println()

	fmt.Println("3. Sharing products...")

	// User1 shares product1 with user2
	err = service.ShareProduct(product1.ID, user1, user2)
	if err != nil {
		log.Printf("Failed to share product: %v", err)
	} else {
		fmt.Printf("Product %s has been shared with user %s\n", product1.Name, user2)
	}

	// User1 shares product2 with user3
	err = service.ShareProduct(product2.ID, user1, user3)
	if err != nil {
		log.Printf("Failed to share product: %v", err)
	} else {
		fmt.Printf("Product %s has been shared with user %s\n", product2.Name, user3)
	}

	fmt.Println()

	fmt.Println("4. Testing permissions after sharing...")

	// Test user2 accessing shared product
	fmt.Printf("User %s accessing shared product %s (read): %t\n",
		user2, product1.Name, service.CanAccessProduct(user2, product1.ID, "read"))

	fmt.Printf("User %s accessing shared product %s (write): %t\n",
		user2, product1.Name, service.CanAccessProduct(user2, product1.ID, "write"))

	// Test user3 accessing shared product
	fmt.Printf("User %s accessing shared product %s (read): %t\n",
		user3, product2.Name, service.CanAccessProduct(user3, product2.ID, "read"))

	fmt.Println()

	fmt.Println("5. Listing user products...")

	// List user1's products
	owned1, shared1 := service.ListUserProducts(user1)
	fmt.Printf("User %s owned products: %d\n", user1, len(owned1))
	for _, p := range owned1 {
		fmt.Printf("  - %s (ID: %s)\n", p.Name, p.ID)
	}
	fmt.Printf("User %s shared products: %d\n", user1, len(shared1))

	// List user2's products
	owned2, shared2 := service.ListUserProducts(user2)
	fmt.Printf("User %s owned products: %d\n", user2, len(owned2))
	for _, p := range owned2 {
		fmt.Printf("  - %s (ID: %s)\n", p.Name, p.ID)
	}
	fmt.Printf("User %s shared products: %d\n", user2, len(shared2))
	for _, p := range shared2 {
		fmt.Printf("  - %s (ID: %s, Owner: %s)\n", p.Name, p.ID, p.Owner)
	}

	fmt.Println()

	fmt.Println("6. Unsharing products...")

	// User1 unshares with user2
	err = service.UnshareProduct(product1.ID, user1, user2)
	if err != nil {
		log.Printf("Failed to unshare product: %v", err)
	} else {
		fmt.Printf("Product %s has been unshared with user %s\n", product1.Name, user2)
	}

	// Test permissions after unsharing
	fmt.Printf("User %s accessing unshared product %s (read): %t\n",
		user2, product1.Name, service.CanAccessProduct(user2, product1.ID, "read"))

	fmt.Println()
	fmt.Println("=== Example completed ===")
}
