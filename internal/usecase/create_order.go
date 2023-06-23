package usecase

import (
	"github.com/santosdvlpr/cleanarq/internal/entity"
	"github.com/santosdvlpr/cleanarq/pkg/events"
)

type (
	OrderInputDTO struct {
		Descricao string `json:"descricao"`
		Preco float64 `json:"preco"`
		Taxa  float64 `json:"taxa"`
	}
	OrderOutputDTO struct {
		ID         int  `json:"id"`
		Descricao  string  `json:"descricao"`
		Preco      float64 `json:"preco"`
		Taxa       float64 `json:"taxa"`
		PrecoTotal float64 `json:"preco_total"`
	}

	CreateOrderUseCase struct {
		OrderRepository entity.OrderRepositoryInterface
		OrderCreated    events.EventInterface
		EventDispatcher events.EventDispatcherInterface
	}
)

func NewCreateOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: OrderRepository,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order := entity.Order{Descricao: input.Descricao, Preco: input.Preco, Taxa: input.Taxa}
	order.CalculaPrecoTotal()
	if err := c.OrderRepository.Save(&order); err != nil {
		return OrderOutputDTO{}, err
	}
	dto := OrderOutputDTO{ID: order.ID, Descricao: order.Descricao, Preco: order.Preco, Taxa: order.Taxa, PrecoTotal: order.PrecoTotal}
	//c.OrderCreated = &events.TestEvent{Name: "OrderCreated", Payload: dto}
	c.OrderCreated.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.OrderCreated)
	return dto, nil
}
