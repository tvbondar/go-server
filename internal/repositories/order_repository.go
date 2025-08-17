package repositories

import "internal/entities"

type OrderRepository interface {
	SaveOrder(order entities.Order) error
	GetOrderByID(id string) (entities.Order, error)
}
