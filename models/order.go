package models

type Order struct {
	ID     int          `json:"id"`
	UserID int          `json:"user_id"`
	Status string       `json:"status"`
	Items  []OrderItems `json:"items"`
}

type OrderItems struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	UnitPrice float64 `json:"unit_price"`
	Quantity  int     `json:"quantity"`
}
