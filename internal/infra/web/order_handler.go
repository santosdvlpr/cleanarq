package web

import (
	"encoding/json"
	"net/http"

	"github.com/santosdvlpr/cleanarq/internal/entity"
	"github.com/santosdvlpr/cleanarq/internal/usecase"
	"github.com/santosdvlpr/cleanarq/pkg/events"
)

type (
	WebOrderHandler struct {
		EventDispatcher   events.EventDispatcherInterface
		OrderRepository   entity.OrderRepositoryInterface
		OrderCreatedEvent events.EventInterface
	}
)

func NewWebOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreatedEvent events.EventInterface,
) *WebOrderHandler {
	return &WebOrderHandler{
		EventDispatcher:   EventDispatcher,
		OrderRepository:   OrderRepository,
		OrderCreatedEvent: OrderCreatedEvent,
	}
}

func (h *WebOrderHandler) List(w http.ResponseWriter, r *http.Request) {

	listOrder := usecase.NewListOrderUseCase(h.OrderRepository)
	output := listOrder.Execute()
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}


func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

