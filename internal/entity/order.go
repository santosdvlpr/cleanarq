package entity

import "errors"

type Order struct {
	ID         int
	Descricao  string
	Preco      float64
	Taxa       float64
	PrecoTotal float64
}

func NewOrder(descricao string, preco float64, taxa float64) (*Order, error) {
	order := &Order{Descricao: descricao, Preco: preco, Taxa: taxa, PrecoTotal: 0}
	err := order.IsValid()
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *Order) IsValid() error {
	if o.Descricao == "" {
		return errors.New("descricao inválida")
	}
	if o.Preco <= 0 {
		return errors.New("preco inválido")
	}
	if o.Taxa <= 0 {
		return errors.New("taxa inválida")
	}
	return nil
}
func (o *Order) CalculaPrecoTotal() error {
	o.PrecoTotal = o.Preco + o.Taxa
	err := o.IsValid()
	if err != nil {
		return err
	}
	return nil
}

type Queries struct {
	db DBTX
}
