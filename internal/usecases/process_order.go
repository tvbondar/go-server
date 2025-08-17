package usecases

import (
	"encoding/json"
	"internal/entities"
	"internal/repositories"
)

type ProcessOrderUseCase struct {
	repo repositories.OrderRepository
}

func NewProcessOrderUseCase(repo repositories.OrderRepository) *ProcessOrderUseCase {
	return &ProcessOrderUseCase{repo: repo}
}

func (u *ProcessOrderUseCase) Execute(rawMessage []byte) error {
	var order entities.Order
	if err := json.Unmarshal(rawMessage, &order); err != nil {
		return err
	}
	// Валидация бизнес-логики
	return u.repo.SaveOrder(order)
}
