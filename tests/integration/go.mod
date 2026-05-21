module github.com/shopee-clone/shopee/tests/integration

go 1.26.3

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/shopee-clone/shopee/packages/go-shared v0.0.0-00010101000000-000000000000
	github.com/shopee-clone/shopee/services/cart v0.0.0-00010101000000-000000000000
	github.com/shopee-clone/shopee/services/inventory v0.0.0-00010101000000-000000000000
	github.com/shopee-clone/shopee/services/order v0.0.0-00010101000000-000000000000
	github.com/shopee-clone/shopee/services/payment v0.0.0-00010101000000-000000000000
	github.com/shopee-clone/shopee/services/promotion v0.0.0-00010101000000-000000000000
)

replace (
	github.com/shopee-clone/shopee/packages/go-shared => ../../packages/go-shared
	github.com/shopee-clone/shopee/services/cart => ../../services/cart
	github.com/shopee-clone/shopee/services/inventory => ../../services/inventory
	github.com/shopee-clone/shopee/services/order => ../../services/order
	github.com/shopee-clone/shopee/services/payment => ../../services/payment
	github.com/shopee-clone/shopee/services/promotion => ../../services/promotion
)
