package repositories

import (
	"database/sql"
	"ecommerce/models"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (repo *OrderRepository) CreateOrder(userID int, order *models.Order) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}

	err = tx.QueryRow("INSERT INTO pedidos(user_id) VALUES ($1) RETURNING id", userID).Scan(&order.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		var price float64
		var stock int
		err := tx.QueryRow("SELECT price, quantity FROM produtos WHERE id = $1", item.ProductID).Scan(&price, &stock)
		if err != nil {
			tx.Rollback()
			return err
		}

		if item.Quantity > stock {
			tx.Rollback()
			return sql.ErrNoRows
		}

		_, err = tx.Exec("UPDATE produtos SET quantity = quantity - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec("INSERT INTO itens_pedido(order_id, product_id, quantity, unit_price) VALUES ($1, $2, $3, $4)",
			order.ID, item.ProductID, item.Quantity, price)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo *OrderRepository) ListOrdersUser(userID int) ([]models.Order, error) {
	rows, err := repo.DB.Query("SELECT id, status FROM pedidos WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		o.UserID = userID
		rows.Scan(&o.ID, &o.Status)
		items, _ := repo.DB.Query("SELECT id, order_id, product_id, quantity, unit_price FROM itens_pedido WHERE order_id = $1", o.ID)
		for items.Next() {
			var item models.OrderItems
			items.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice)
			o.Items = append(o.Items, item)
		}

		orders = append(orders, o)
	}
	return orders, nil
}

func (repo *OrderRepository) OrderPayment(userID int, orderID int) error {
	var status string
	err := repo.DB.QueryRow("SELECT status FROM pedidos WHERE id = $1 AND user_id = $2", orderID, userID).Scan(&status)
	if err != nil {
		return err
	}

	if status == "pago" {
		return sql.ErrNoRows
	}

	_, err = repo.DB.Exec("UPDATE pedidos SET status = 'pago' WHERE id = $1", orderID)
	return err
}

func (repo *OrderRepository) ListAllOrdersAdmin() ([]models.Order, error) {
	rows, err := repo.DB.Query("SELECT id, user_id, status FROM pedidos ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status); err != nil {
			continue
		}

		items, _ := repo.DB.Query("SELECT product_id, quantity, unit_price FROM itens_pedido WHERE order_id = $1", o.ID)
		for items.Next() {
			var item models.OrderItems
			items.Scan(&item.ProductID, &item.Quantity, &item.UnitPrice)
			o.Items = append(o.Items, item)
		}

		orders = append(orders, o)
	}
	return orders, nil
}

func (repo *OrderRepository) CancelOrder(userID int, orderID int) error {
	var status string
	err := repo.DB.QueryRow("SELECT status FROM pedidos WHERE id = $1 AND user_id = $2", orderID, userID).Scan(&status)
	if err != nil {
		return err
	}

	if status != "pendente" {
		return sql.ErrNoRows
	}

	rows, err := repo.DB.Query("SELECT product_id, quantity FROM itens_pedido WHERE order_id = $1", orderID)
	if err != nil {
		return err
	}
	defer rows.Close()

	type item struct {
		ProductID int
		Quantity  int
	}
	var items []item
	for rows.Next() {
		var i item
		rows.Scan(&i.ProductID, &i.Quantity)
		items = append(items, i)
	}

	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}

	for _, i := range items {
		_, err := tx.Exec("UPDATE produtos SET quantity = quantity + $1 WHERE id = $2", i.Quantity, i.ProductID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.Exec("UPDATE pedidos SET status = 'cancelado' WHERE id = $1", orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo *OrderRepository) GetOrderCounts() (models.CountOrders, error) {
	var count models.CountOrders
	var querie = "SELECT COUNT(CASE WHEN status = 'cancelado' THEN 1 END) AS cancelados, COUNT(CASE WHEN status = 'pendente' THEN 1 END) AS pendentes, COUNT(CASE WHEN status = 'pago' THEN 1 END) AS pagos, (SELECT SUM(ip.quantity * ip.unit_price) FROM itens_pedido ip JOIN pedidos p2 ON p2.id = ip.order_id WHERE p2.status = 'pago') AS total_valor_pago  FROM pedidos p"
	err := repo.DB.QueryRow(querie).Scan(&count.Canceled, &count.Pending, &count.Paid, &count.TotalPaid)
	if err != nil {
		return count, err
	}
	return count, nil
}
