package usecases

import (
	"encoding/json"
	"fmt"

	"github.com/tvbondar/go-server/internal/repositories"

	"github.com/tvbondar/go-server/internal/entities"
)

type ProcessOrderUseCase struct {
	dbRepo    repositories.OrderRepository
	cacheRepo repositories.CacheOrderRepository
}

func NewProcessOrderUseCase(dbRepo repositories.OrderRepository, cacheRepo repositories.CacheOrderRepository) *ProcessOrderUseCase {
	return &ProcessOrderUseCase{dbRepo: dbRepo, cacheRepo: cacheRepo}
}

func (u *ProcessOrderUseCase) Execute(rawMessage []byte) error {
	var order entities.Order
	if err := json.Unmarshal(rawMessage, &order); err != nil {
		fmt.Println("Invalid message:", err) // Логируем и игнорируем
		return nil                           // Не ошибка, просто игнор
	}
	// Валидация (например, order.OrderUID != "")
	if order.OrderUID == "" {
		fmt.Println("Invalid order UID")
		return nil
	}
	// Сохраняем в DB и кэш
	if err := u.dbRepo.SaveOrder(order); err != nil {
		return err
	}
	u.cacheRepo.SaveOrder(order)
	return nil
}
