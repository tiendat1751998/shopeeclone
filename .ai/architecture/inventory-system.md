# Concurrency Inventory System - Flash Sale

During high-concurrency campaigns (Flash Sales), millions of requests try to buy limited stock. We intercept this bottleneck using **Redis Memory Stock Deduct via Lua Scripts**.

## Production Lua Script: Stock Reservation (`reserve.lua`)
This script must be loaded into Redis via `SCRIPT LOAD`. The service executes it using `EVALSHA`.

```lua
-- KEYS[1]: stock key (e.g. "stock:sku_12345")
-- KEYS[2]: user order lock key (e.g. "order_lock:user_12:sku_12345")
-- ARGV[1]: order quantity (e.g. 1)
-- ARGV[2]: dynamic expiration time of lock (in seconds, e.g. 300)

local stock_key = KEYS[1]
local lock_key = KEYS[2]
local qty = tonumber(ARGV[1])
local lock_ttl = tonumber(ARGV[2])

-- 1. Check if user already holds a pending purchase lock
if redis.call("EXISTS", lock_key) == 1 then
    return -1 -- Error: User has pending reservation
end

-- 2. Fetch current in-memory stock level
local current_stock = redis.call("GET", stock_key)
if not current_stock then
    return -2 -- Error: Stock key does not exist
end

current_stock = tonumber(current_stock)
if current_stock < qty then
    return 0 -- Error: Out of stock (Insufficient quantity)
end

-- 3. Atomically deduct stock and record lock
redis.call("DECRBY", stock_key, qty)
redis.call("SET", lock_key, qty, "EX")
redis.call("EXPIRE", lock_key, lock_ttl)

return 1 -- Success: Stock reserved!
```

## Compensation Lua Script: Release Reservation (`release.lua`)
Used if payment fails, canceling the transaction.
```lua
local stock_key = KEYS[1]
local lock_key = KEYS[2]

local reserved_qty = redis.call("GET", lock_key)
if reserved_qty then
    redis.call("INCRBY", stock_key, tonumber(reserved_qty))
    redis.call("DEL", lock_key)
    return 1 -- Success: Stock returned!
else
    return 0 -- No reservation to release
end
```
