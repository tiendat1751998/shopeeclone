package http

import (
	"github.com/shopee-clone/shopee/platforms/oms-fulfillment/internal/inventory"
	"github.com/shopee-clone/shopee/platforms/oms-fulfillment/internal/ordermanagement"
	"github.com/shopee-clone/shopee/platforms/oms-fulfillment/internal/pickpack"
	"github.com/shopee-clone/shopee/platforms/oms-fulfillment/internal/returns"
	"github.com/shopee-clone/shopee/platforms/oms-fulfillment/internal/warehouse"
)

type Handler struct {
	orders    *ordermanagement.Service
	inventory *inventory.Service
	pickpack  *pickpack.Service
	returns   *returns.Service
	warehouse *warehouse.Service
}

func NewHandler(
	orders *ordermanagement.Service,
	inventory *inventory.Service,
	pickpack *pickpack.Service,
	returns *returns.Service,
	warehouse *warehouse.Service,
) *Handler {
	return &Handler{
		orders:    orders,
		inventory: inventory,
		pickpack:  pickpack,
		returns:   returns,
		warehouse: warehouse,
	}
}
