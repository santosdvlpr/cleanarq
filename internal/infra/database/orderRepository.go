package database

import (
	"database/sql"
	"fmt"

	"github.com/santosdvlpr/cleanarq/internal/entity"
)

type (
	OrderRepository struct {
		Db *sql.DB
	}
)

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (o *OrderRepository) Save(order *entity.Order) error {
	stmt, err := o.Db.Prepare("insert into orders(descricao,preco,taxa,preco_total) values(?,?,?,?);")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(&order.Descricao, &order.Preco, &order.Taxa, &order.PrecoTotal)
	if err != nil {
		return err
	}
	return nil
}

const orderList = `-- name: OrderList :many
SELECT id, descricao, preco, taxa, preco_total FROM orders
WHERE descricao <> '' ORDER BY descricao 
`

func (o *OrderRepository) List() []entity.Order {
	rows, err := o.Db.Query(orderList)
	if err != nil {
		fmt.Println("deu erro aqui...")
		panic(err)
	}

	defer rows.Close()
	var orders []entity.Order
	var order entity.Order
	for rows.Next() {
		if err := rows.Scan(
			&order.ID, &order.Descricao, &order.Preco, &order.Taxa, &order.PrecoTotal,
		); err != nil {
			fmt.Println("deu erro aqui........")
			panic(err)
		}
		fmt.Printf("id:%v descrição:%s\n", order.ID, order.Descricao)
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("tamanho do slice:%v\n", len(orders))
	fmt.Printf("orders: %v\n", orders)
	return orders
}

