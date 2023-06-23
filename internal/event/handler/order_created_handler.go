package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/santosdvlpr/cleanarq/pkg/events"
)

type (
	OrderCreatedHandler struct {
		RabbitMQChannel *amqp.Channel
	}
)

func NewOrderCreatedHandler(RabbitMQChannel *amqp.Channel) *OrderCreatedHandler {
	return &OrderCreatedHandler{RabbitMQChannel}
}
func (h *OrderCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Ordem de serviço criada:%v", event.GetPayload())
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}
	h.RabbitMQChannel.Publish(
		"amq.direct", // exchange
		"",           //key name
		false,        //madatory
		false,        //immediate
		msgRabbitmq,  //message to publish
	)
}
