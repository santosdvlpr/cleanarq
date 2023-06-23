package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"context"
	"fmt"

	"github.com/santosdvlpr/cleanarq/internal/infra/graph/model"
	"github.com/santosdvlpr/cleanarq/internal/usecase"
)

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, input *model.OrderInput) (*model.Order, error) {
	dto := usecase.OrderInputDTO{Descricao: input.Descricao, Preco: input.Preco, Taxa: input.Taxa}
	output, err := r.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &model.Order{
		ID:         output.ID,
		Descricao:  output.Descricao,
		Preco:      output.Preco,
		Taxa:       output.Taxa,
		PrecoTotal: output.PrecoTotal,
	}, nil
}

// ListOrders is the resolver for the ListOrders field.
func (r *queryResolver) ListOrders(ctx context.Context) ([]*model.Order, error) {
	print("ENTRANDO....")
	list := r.ListOrderUseCase.OrderRepository.List()
	print("SAINDO....")
	var orders []*model.Order
	for _, order := range list {
		fmt.Println("Order:", order.Descricao)
		orders = append(orders, &model.Order{
			ID:         order.ID,
			Descricao:  order.Descricao,
			Preco:      order.Preco,
			Taxa:       order.Taxa,
			PrecoTotal: order.PrecoTotal,
		})
	}
	print("retornando....:", orders)
	return orders, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
