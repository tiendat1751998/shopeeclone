# Redis Caching & Memory Strategies

Redis is the caching backbone of our high-scale system.

## 1. Cache Patterns

### Cache-Aside Code Template (Golang)
```go
func GetProductCatalog(ctx context.Context, spuID string) (*Product, error) {
    cacheKey := fmt.Sprintf("product:%s", spuID)
    
    // 1. Read from Redis
    if val, err := rdb.Get(ctx, cacheKey).Result(); err == nil {
        var product Product
        json.Unmarshal([]byte(val), &product)
        return &product, nil
    }
    
    // 2. Cache Miss: Read from MongoDB source of truth
    product, err := mongoDb.FetchSPU(ctx, spuID)
    if err != nil {
        return nil, err
    }
    
    // 3. Write back to cache (TTL: 1 hour)
    data, _ := json.Marshal(product)
    rdb.Set(ctx, cacheKey, data, 1 * time.Hour)
    
    return product, nil
}
```

## 2. Threat Mitigations
- **Cache Stampede Prevention**: Use **Single-Flight** execution locks so that during a cache miss under high concurrency, only **one** thread queries the database, while others wait for its return.
- **Cache Penetration Protection (Bloom Filters)**: Deploy Redis Bloom Filters (`BF.EXISTS`, `BF.ADD`) for non-existent SPU/SKU IDs. If the ID is not in the Bloom Filter, reject the HTTP request immediately without hitting the MongoDB.
