//go:build wireinject
// +build wireinject

package main

// o wire vai processar este código e gerar o arquivo wire_gen.go

import (
	"database/sql"
	"github.com/google/wire"
	"github.com/santosdvlpr/cleanarq/internal/entity"
	"github.com/santosdvlpr/cleanarq/internal/event"
	"github.com/santosdvlpr/cleanarq/internal/infra/database"
	"github.com/santosdvlpr/cleanarq/internal/infra/web"
	"github.com/santosdvlpr/cleanarq/internal/usecase"
	"github.com/santosdvlpr/cleanarq/pkg/events"
)

// Cria um Set (p/o OrderRepository) passando sua interface e sua implementação concreta para o Binde do wire
// A saida setOrderRepositoryDependency representa o repositório com todas as suas dependências injetadas
// pelo wire
var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

// Cria um Set (p/o EventDispatcher) passando sua interface e sua implementação concreta para o Binde do wire
// A saida "setEventDispatcherDependency" representa o Dispatcher com todas as suas dependências injetadas
// pelo wire
var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

// Cria um Set (p/o OrderCreatedEvent) passando sua interface e sua implementação concreta para o Binde do wire
// A saida setOrderCreatedEvent representa o evento com todas as suas dependências injetadas pelo wire
var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

// Bilda (constroi) um CreateOrderUseCase passando os SETs com todas as dependências injetadas e mais: conexão de banco de dados e a interface do eventDispataher
func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

// Bilda (constroi) um ListOrderUseCase passando os SETs com todas as dependências injetadas e mais: conexão de banco de dados
func NewListOrderUseCase(db *sql.DB) *usecase.ListOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewListOrderUseCase,
	)
	return &usecase.ListOrderUseCase{}
}

// Bilda (constroi) um WebOrderHandler passando os SETs com todas as dependências injetadas e mais: conexão de banco de dados e a interface do eventDispataher  --- >>  caso se esteja trabalhando com rest api endpoit
func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
