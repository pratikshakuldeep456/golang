package fooddelivery

type PaymentStatus string

const (
	PENDING PaymentStatus = "pending"
	PAID    PaymentStatus = "paid"
	FAILED  PaymentStatus = "failed"
)

type OrderStatus string

const (
	ORDERPENDING OrderStatus = "ongoing"
	ORDERDELIVED OrderStatus = "delievered"
	ORDERFAILED  OrderStatus = "failed"
)
