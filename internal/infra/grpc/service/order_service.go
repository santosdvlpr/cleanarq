package service

import (
	"context"

	"github.com/santosdvlpr/cleanarq/internal/infra/grpc/pb"
	"github.com/santosdvlpr/cleanarq/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrderUseCase   usecase.ListOrderUseCase
}

func NewOrderService(
	createOrderUseCase usecase.CreateOrderUseCase,
	listOrderUseCase usecase.ListOrderUseCase,
) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		Descricao: string(in.Descricao),
		Preco:     float64(in.Preco),
		Taxa:      float64(in.Taxa),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         int32(output.ID),
		Descricao:  string(output.Descricao),
		Preco:      float32(output.Preco),
		Taxa:       float32(output.Taxa),
		PrecoTotal: float32(output.PrecoTotal),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.Blank) (*pb.OrderList, error) {
	orders := s.ListOrderUseCase.Execute()
	var list []*pb.Order
	for _, orderDTO := range orders {
		order := &pb.Order{
			Id:         int32(orderDTO.ID),
			Descricao:  orderDTO.Descricao,
			Preco:      float32(orderDTO.Preco),
			Taxa:       float32(orderDTO.Taxa),
			PrecoTotal: float32(orderDTO.PrecoTotal),
		}
		list = append(list, order)
	}
	return &pb.OrderList{Orders: list}, nil
}
