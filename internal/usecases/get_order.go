package usecases

import (
	"github.com/tvbondar/go-server/internal/entities"
	"github.com/tvbondar/go-server/internal/repositories"
)

type GetOrderUseCase struct {
	cacheRepo repositories.OrderRepository
	dbRepo    repositories.OrderRepository
}

func NewGetOrderUseCase(cacheRepo, dbRepo repositories.OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{cacheRepo: cacheRepo, dbRepo: dbRepo}
}

func (u *GetOrderUseCase) Execute(id string) (entities.Order, error) {
	order, err := u.cacheRepo.GetOrderByID(id)
	if err == nil {
		return order, nil
	}
	order, err = u.dbRepo.GetOrderByID(id)
	if err != nil {
		return entities.Order{}, err
	}
	u.cacheRepo.SaveOrder(order)
	return order, nil
}
