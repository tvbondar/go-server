package repositories

import (
	"database/sql"

	"github.com/tvbondar/go-server/internal/entities"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) SaveOrder(order entities.Order) error {
	// Используй транзакцию для атомарности
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Вставка в orders
	_, err = tx.Exec("INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}

	// Вставка в delivery
	_, err = tx.Exec("INSERT INTO deliveries (name, phone, zip, city, address, region, email, order_uid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email, order.OrderUID)
	if err != nil {
		return err
	}

	// Аналогично для payment и items (добавь циклы для items)
	// Для items:
	for _, item := range order.Items {
		_, err = tx.Exec("INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_uid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status, order.OrderUID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresOrderRepository) GetOrderByID(id string) (entities.Order, error) {
	var order entities.Order
	// Запрос SELECT из orders, deliveries, payments, items
	// Пример для orders
	err := r.db.QueryRow("SELECT * FROM orders WHERE order_uid = $1", id).Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		return order, err
	}
	// Аналогично для delivery, payment, items (используй Query для items)
	rows, err := r.db.Query("SELECT * FROM items WHERE order_uid = $1", id)
	if err != nil {
		return order, err
	}
	for rows.Next() {
		var item entities.Item
		rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		order.Items = append(order.Items, item)
	}
	// Добавь для delivery и payment
	return order, nil
}

func (r *PostgresOrderRepository) GetAllOrders() ([]entities.Order, error) {
	var orders []entities.Order
	// SELECT * FROM orders, затем для каждого order_uid загрузи delivery, payment, items
	rows, err := r.db.Query("SELECT order_uid FROM orders")
	if err != nil {
		return orders, err
	}
	for rows.Next() {
		var id string
		rows.Scan(&id)
		order, err := r.GetOrderByID(id)
		if err != nil {
			continue // Или обработай ошибку
		}
		orders = append(orders, order)
	}
	return orders, nil
}
