package repositories

import (
	"fmt"

	"github.com/tvbondar/go-server/internal/entities"
)

type CacheOrderRepository struct {
	cache map[string]entities.Order
}

func NewCacheOrderRepository() *CacheOrderRepository {
	return &CacheOrderRepository{cache: make(map[string]entities.Order)}
}

func (r *CacheOrderRepository) SaveOrder(order entities.Order) error {
	r.cache[order.OrderUID] = order
	return nil
}

func (r *CacheOrderRepository) GetOrderByID(id string) (entities.Order, error) {
	order, exists := r.cache[id]
	if !exists {
		return entities.Order{}, fmt.Errorf("order not found")
	}
	return order, nil
}

func (r *CacheOrderRepository) LoadFromDB(dbRepo OrderRepository) error {
	orders, err := dbRepo.GetAllOrders()
	if err != nil {
		return err
	}
	for _, order := range orders {
		r.cache[order.OrderUID] = order
	}
	return nil
}
