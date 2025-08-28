# Casbin Product Sharing Example

This example demonstrates how to use Casbin to implement Attribute-Based Access Control (ABAC) for managing product sharing functionality.


## Architecture Design

### Why Choose ABAC over RBAC?

**ABAC Advantages**: Dynamic permissions, attribute-based, flexible expansion. 
**RBAC Problems**: Requires predefined roles, complex maintenance, poor scalability.

### Permission Model

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (r.sub == r.obj.Owner) || contains(r.obj.SharedUsers, r.sub)
```

**Permission Rule Explanation**:
- `r.sub == r.obj.Owner`: Product owner has full access permissions
- `r.sub in r.obj.SharedUsers`: Shared users have access permissions
- Supported operations: `read`, `write`, etc.

## File Structure

```
product_sharing_example/
├── model.conf              # Casbin model configuration
├── product.go              # Product struct and methods
├── product_service.go      # Product service logic
├── main.go                 # Main program example
├── product_service_test.go # Test files
└── README.md               # Documentation
```

## Usage

### 1. Run Example

```bash
cd product_sharing_example
go run .
```

### 2. Run Tests

```bash
go test -v
```

### 3. Integrate into Your Project

```go
// Create product service
service, err := NewProductService()
if err != nil {
    log.Fatal(err)
}
defer service.Close()

// Create product
product, err := service.CreateProduct("prod_001", "Product Name", "user@example.com")
if err != nil {
    log.Fatal(err)
}

// Share product
err = service.ShareProduct("prod_001", "user@example.com", "other@example.com")
if err != nil {
    log.Fatal(err)
}

// Check permissions
canAccess := service.CanAccessProduct("other@example.com", "prod_001", "read")
```

## API Interface

### ProductService

- `CreateProduct(id, name, ownerID)`: Create new product
- `GetProduct(id)`: Get product information
- `ShareProduct(productID, ownerID, targetUserID)`: Share product
- `UnshareProduct(productID, ownerID, targetUserID)`: Unshare product
- `CanAccessProduct(userID, productID, action)`: Check access permissions
- `ListUserProducts(userID)`: List user's products

### Product

- `AddSharedUser(userID)`: Add shared user
- `RemoveSharedUser(userID)`: Remove shared user
- `IsSharedUser(userID)`: Check if user is a shared user
- `GetSharedUsers()`: Get all shared users

## Dependencies

- Go 1.16+
- github.com/casbin/casbin/v2
