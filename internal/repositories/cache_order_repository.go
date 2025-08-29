package repositories

import (
	"fmt"
	"sync"

	"github.com/tvbondar/go-server/internal/entities"
)

type CacheOrderRepository struct {
	cache map[string]entities.Order
	mu    sync.RWMutex
}

func NewCacheOrderRepository() *CacheOrderRepository {
	return &CacheOrderRepository{
		cache: make(map[string]entities.Order),
	}
}

func (r *CacheOrderRepository) SaveOrder(order entities.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[order.OrderUID] = order
	return nil
}

func (r *CacheOrderRepository) GetOrderByID(id string) (entities.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, exists := r.cache[id]
	if !exists {
		return entities.Order{}, fmt.Errorf("order not found")
	}
	return order, nil
}

func (r *CacheOrderRepository) GetAllOrders() ([]entities.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var orders []entities.Order
	for _, order := range r.cache {
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *CacheOrderRepository) LoadFromDB(dbRepo OrderRepository) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	orders, err := dbRepo.GetAllOrders()
	if err != nil {
		return err
	}
	for _, order := range orders {
		r.cache[order.OrderUID] = order
	}
	fmt.Println("Cache loaded with", len(orders), "orders")
	return nil
}
