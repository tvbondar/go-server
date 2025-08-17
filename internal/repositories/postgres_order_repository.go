package repositories

import (
	"database/sql"
	"internal/entities"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) SaveOrder(order entities.Order) error {
	// SQL-вставка в таблицы orders, deliveries, payments, items
	_, err := r.db.Exec("INSERT INTO orders (order_uid, track_number) VALUES ($1, $2)", order.OrderUID, order.TrackNumber)
	// Аналогично для других таблиц
	return err
}

func (r *PostgresOrderRepository) GetOrderByID(id string) (entities.Order, error) {
	var order entities.Order
	// SQL-запрос для получения данных
	return order, nil
}
