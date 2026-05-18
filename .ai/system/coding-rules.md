# Coding Rules & Best Practices

## 1. Go (Golang) Backend Services
- **Concurrency Management**: Never start a goroutine without an explicit panic recovery wrapper and context propagation. Use `golang.org/x/sync/errgroup` to orchestrate parallel workers.
  ```go
  g, ctx := errgroup.WithContext(parentCtx)
  g.Go(func() error {
      defer func() {
          if r := recover(); r != nil {
              log.Errorf("Panic recovered: %v", r)
          }
      }()
      return processItem(ctx, item)
  })
  ```
- **GC Optimization**: For high-throughput services (Inventory, Cart), minimize memory allocation. Reuse byte slices and structs using `sync.Pool`.
- **Error Handling**: Use structured error wrapping. Always expose user-friendly error codes while logging root-cause system logs:
  ```go
  if err != nil {
      return fmt.Errorf("%w: failed to fetch inventory for SKU %s", ErrInternalDatabase, skuID)
  }
  ```

## 2. Java Spring Boot Services
- **JPA & Hibernate**:
  - All associations must be `FetchType.LAZY` by default to prevent N+1 select queries.
  - Use `@EntityGraph` or custom JPQL joins for queries requiring associated records.
  - Transactions must have explicit boundaries using `@Transactional(propagation = Propagation.REQUIRED, isolation = Isolation.READ_COMMITTED)`.
- **API Validation**: Enforce strict constraints on all incoming DTO requests using Jakarta Validation:
  ```java
  public record CreateOrderRequest(
      @NotNull(message = "Cart items list cannot be null")
      @NotEmpty(message = "Cart items cannot be empty")
      List<@Valid OrderItemDTO> items,
      @NotBlank(message = "Shipping address is required")
      String shippingAddress,
      @NotBlank(message = "Payment method code is required")
      String paymentMethodCode
  ) {}
  ```

## 3. Next.js 15 (TypeScript) Frontend
- **Rendering Strategy**: Leverage Server Components (RSC) for initial page loads (SEO optimized), and Client Components (`"use client"`) only for interactive islands (Checkout buttons, variation selector modal).
- **Type Safety**: Strictly avoid the `any` type. Define explicit interfaces for all API responses and component props.
- **State Hydration**: Keep React components synchronized with Redux Toolkit or Server Actions using dynamic state revalidation: `revalidatePath("/cart")`.
