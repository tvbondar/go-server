package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/tvbondar/go-server/internal/entities"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) SaveOrder(order entities.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	deliveryJSON, err := json.Marshal(order.Delivery)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery: %w", err)
	}
	paymentJSON, err := json.Marshal(order.Payment)
	if err != nil {
		return fmt.Errorf("failed to marshal payment: %w", err)
	}

	_, err = tx.Exec(`
        INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, delivery, payment)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
		deliveryJSON, paymentJSON)
	if err != nil {
		return fmt.Errorf("failed to insert order %s: %w", order.OrderUID, err)
	}

	for _, item := range order.Items {
		_, err = tx.Exec(`
            INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("failed to insert item for order %s: %w", order.OrderUID, err)
		}
	}

	return tx.Commit()
}

func (r *PostgresOrderRepository) GetOrderByID(id string) (entities.Order, error) {
	var order entities.Order
	var deliveryJSON, paymentJSON []byte

	err := r.db.QueryRow(`
        SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, delivery, payment
        FROM orders WHERE order_uid = $1`,
		id).Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
		&deliveryJSON, &paymentJSON)
	if err != nil {
		return entities.Order{}, fmt.Errorf("failed to query order %s: %w", id, err)
	}

	if err := json.Unmarshal(deliveryJSON, &order.Delivery); err != nil {
		return entities.Order{}, fmt.Errorf("failed to unmarshal delivery: %w", err)
	}
	if err := json.Unmarshal(paymentJSON, &order.Payment); err != nil {
		return entities.Order{}, fmt.Errorf("failed to unmarshal payment: %w", err)
	}

	rows, err := r.db.Query(`
        SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
        FROM items WHERE order_uid = $1`, id)
	if err != nil {
		return entities.Order{}, fmt.Errorf("failed to query items for order %s: %w", id, err)
	}
	defer rows.Close()
	for rows.Next() {
		var item entities.Item
		err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return entities.Order{}, fmt.Errorf("failed to scan item for order %s: %w", id, err)
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *PostgresOrderRepository) GetAllOrders() ([]entities.Order, error) {
	var orders []entities.Order
	rows, err := r.db.Query("SELECT order_uid FROM orders")
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			continue
		}
		order, err := r.GetOrderByID(id)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}
	return orders, nil
}
