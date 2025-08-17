package repositories

import "github.com/tvbondar/go-server/internal/entities"

type OrderRepository interface {
	SaveOrder(order entities.Order) error
	GetOrderByID(id string) (entities.Order, error)
	GetAllOrders() ([]entities.Order, error) // Для восстановления кэша
}
