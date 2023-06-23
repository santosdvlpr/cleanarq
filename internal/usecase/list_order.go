package usecase

import (
	"github.com/santosdvlpr/cleanarq/internal/entity"
)

type (
	ListOrderUseCase struct {
		OrderRepository entity.OrderRepositoryInterface
	}
)

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

func (l *ListOrderUseCase) Execute() []*OrderOutputDTO {
	var list []*OrderOutputDTO
	orders := l.OrderRepository.List()
	for _, order := range orders {
		dto := OrderOutputDTO{
			ID: order.ID, Descricao: order.Descricao, Preco: order.Preco, Taxa: order.Taxa, PrecoTotal: order.PrecoTotal,
		}
		list = append(list, &dto)
		//c.OrderCreated = &events.TestEvent{Name: "OrderCreated", Payload: dto}
		//l.OrderCreated.SetPayload(dto)
		//l.EventDispatcher.Dispatch(c.OrderCreated)

	}
	return list
}
