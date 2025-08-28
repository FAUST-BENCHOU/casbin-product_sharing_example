# Casbin 产品共享示例

这个示例演示了如何使用 Casbin 实现基于属性的访问控制（ABAC）来管理产品共享功能。

## 功能特性

- **产品创建**: 用户可以创建产品并成为所有者
- **产品共享**: 产品所有者可以将产品共享给其他用户
- **权限控制**: 基于产品所有者和共享用户列表进行访问控制
- **动态权限**: 权限可以实时添加和移除，无需重启服务

## 架构设计

### 为什么选择 ABAC 而不是 RBAC？

1. **动态权限**: 产品共享关系是动态的，不是预定义的角色
2. **基于属性**: 权限基于产品的 `Owner` 和 `SharedUsers` 属性
3. **灵活性**: 可以轻松添加新的共享用户，无需修改策略配置

### 权限模型

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (r.sub == r.obj.Owner) || (r.sub in r.obj.SharedUsers)
```

**权限规则解释**:
- `r.sub == r.obj.Owner`: 产品所有者拥有完全访问权限
- `r.sub in r.obj.SharedUsers`: 共享用户拥有访问权限
- 支持的操作: `read`, `write` 等

## 文件结构

```
product_sharing_example/
├── model.conf              # Casbin 模型配置
├── product.go              # 产品结构体和方法
├── product_service.go      # 产品服务逻辑
├── main.go                 # 主程序示例
├── product_service_test.go # 测试文件
└── README.md               # 说明文档
```

## 使用方法

### 1. 运行示例

```bash
cd product_sharing_example
go run .
```

### 2. 运行测试

```bash
go test -v
```

### 3. 集成到你的项目

```go
// 创建产品服务
service, err := NewProductService()
if err != nil {
    log.Fatal(err)
}
defer service.Close()

// 创建产品
product, err := service.CreateProduct("prod_001", "产品名称", "user@example.com")
if err != nil {
    log.Fatal(err)
}

// 共享产品
err = service.ShareProduct("prod_001", "user@example.com", "other@example.com")
if err != nil {
    log.Fatal(err)
}

// 检查权限
canAccess := service.CanAccessProduct("other@example.com", "prod_001", "read")
```

## API 接口

### ProductService

- `CreateProduct(id, name, ownerID)`: 创建新产品
- `GetProduct(id)`: 获取产品信息
- `ShareProduct(productID, ownerID, targetUserID)`: 共享产品
- `UnshareProduct(productID, ownerID, targetUserID)`: 取消共享
- `CanAccessProduct(userID, productID, action)`: 检查访问权限
- `ListUserProducts(userID)`: 列出用户的产品

### Product

- `AddSharedUser(userID)`: 添加共享用户
- `RemoveSharedUser(userID)`: 移除共享用户
- `IsSharedUser(userID)`: 检查是否是共享用户
- `GetSharedUsers()`: 获取所有共享用户

## 扩展建议

1. **权限级别**: 可以添加不同的权限级别（如只读、读写、管理等）
2. **时间限制**: 可以添加共享的过期时间
3. **审计日志**: 记录所有的权限变更操作
4. **批量操作**: 支持批量共享和取消共享
5. **通知系统**: 当产品被共享时通知相关用户

## 注意事项

1. **并发安全**: 使用互斥锁保护共享数据
2. **错误处理**: 所有操作都有适当的错误处理
3. **内存管理**: 产品数据存储在内存中，生产环境建议使用数据库
4. **性能**: 对于大量产品，可以考虑使用缓存优化权限检查

## 依赖

- Go 1.16+
- github.com/casbin/casbin/v2
