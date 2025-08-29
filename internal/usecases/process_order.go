package usecases

import (
	"encoding/json"
	"fmt"

	"github.com/tvbondar/go-server/internal/entities"
	"github.com/tvbondar/go-server/internal/repositories"
)

type ProcessOrderUseCase struct {
	dbRepo    repositories.OrderRepository
	cacheRepo repositories.OrderRepository
}

func NewProcessOrderUseCase(dbRepo, cacheRepo repositories.OrderRepository) *ProcessOrderUseCase {
	return &ProcessOrderUseCase{dbRepo: dbRepo, cacheRepo: cacheRepo}
}

func (u *ProcessOrderUseCase) Execute(rawMessage []byte) error {
	var order entities.Order
	if err := json.Unmarshal(rawMessage, &order); err != nil {
		fmt.Printf("Invalid JSON message: %v\n", err)
		return nil
	}
	if order.OrderUID == "" {
		fmt.Println("Invalid order: empty OrderUID")
		return nil
	}
	if order.TrackNumber == "" {
		fmt.Println("Invalid order: empty TrackNumber")
		return nil
	}
	if len(order.Items) == 0 {
		fmt.Println("Invalid order: no items")
		return nil
	}
	if order.Delivery.Name == "" || order.Delivery.Phone == "" {
		fmt.Println("Invalid order: incomplete delivery info")
		return nil
	}
	if order.Payment.Transaction == "" || order.Payment.Amount <= 0 {
		fmt.Println("Invalid order: incomplete payment info")
		return nil
	}

	if err := u.dbRepo.SaveOrder(order); err != nil {
		return err
	}
	return u.cacheRepo.SaveOrder(order)
}
