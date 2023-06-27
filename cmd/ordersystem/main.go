package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/santosdvlpr/cleanarq/configs"
	"github.com/santosdvlpr/cleanarq/internal/event/handler"
	"github.com/santosdvlpr/cleanarq/internal/infra/graph"
	"github.com/santosdvlpr/cleanarq/internal/infra/grpc/pb"
	"github.com/santosdvlpr/cleanarq/internal/infra/grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/santosdvlpr/cleanarq/internal/infra/web/webserver"
	"github.com/santosdvlpr/cleanarq/pkg/events"

	// mysql
	_ "github.com/go-sql-driver/mysql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func init() {
	// Connect to MySQL database
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/orders")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get the driver instance

	//driver, err := mysql.WithInstance(db, &mysql.Config{})
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Create a new migrate instance with the driver and the source of migrations
	m, err := migrate.New(
		"file://../../sqlc/sql/migrations/",            // path to migration files
		"mysql://root:root@tcp(localhost:3306)/orders", // name of the database
		//driver,                                // driver instance
	)
	if err != nil {
		log.Fatal("error aqui:", err)
	}

	// Run up migrations to apply all available migrations

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("executa o up da migração", err)
	}
}

func main() {
	//1) Levantar as configuraçoes
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	//2) Conectar ao banco de dados
	//fmt.Printf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName)
	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//3) Abrir o canal de comunicação
	rabbitMQChannel := getRabbitMQChannel()
	/* 	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	   	if err != nil {
	   		panic(err)
	   	}
	   	rabbitMQChannel, err := conn.Channel()
	   	if err != nil {
	   		panic(err)
	   	}
	*/
	//4) Criar o dispatcher registrando o evento cujo handler recebe o canal de comunicação
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	//5) Instancia o usecase passando banco e o dispatcher
	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrderUseCase(db)

	//6 Subindo os servidores

	// 6.1  web server
	webServer := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webServer.AddHandler("/create", webOrderHandler.Create)
	webServer.AddHandler("/order", webOrderHandler.List)
	fmt.Println("Iniciado web server na porta:", configs.WebServerPort)
	go webServer.Start()

	// 6.2  grpc server
	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)
	fmt.Println("Iniciado gRPC server na porta:", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	// 6.3  graphQL server
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			CreateOrderUseCase: *createOrderUseCase,
			ListOrderUseCase:   *listOrderUseCase,
		},
	}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	fmt.Println("Iniciado servidor GraphQL na porta:", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
