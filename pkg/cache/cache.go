package cache

import (
    "sync"
    "go-server/internal/model"
)

type Cache struct {
    mu     sync.RWMutex
    orders map[string]model.Order
}

func NewCache() *Cache {
    return &Cache{
        orders: make(map[string]model.Order),
    }
}

func (c *Cache) Set(order model.Order) {
    c.mu.Lock()
    c.orders[order.OrderUID] = order
    c.mu.Unlock()
}

func (c *Cache) Get(orderUID string) (model.Order, bool) {
    c.mu.RLock()
    order, exists := c.orders[orderUID]
    c.mu.RUnlock()
    return order, exists
}

func (c *Cache) Load(orders []model.Order) {
    c.mu.Lock()
    for _, order := range orders {
        c.orders[order.OrderUID] = order
    }
    c.mu.Unlock()
}
