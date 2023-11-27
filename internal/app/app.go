package app

import (
	"fmt"
	"log"
	"test_wb/config"
	"test_wb/internal/cache"
	"test_wb/internal/nats"
	"test_wb/internal/orders/controller"
	"test_wb/internal/orders/generator"
	"test_wb/internal/repository"
	"test_wb/pkg/httpserver"
	"test_wb/pkg/postgres"
	"time"
)

func Start(cfg *config.Config) {
	//подключаем натс стриминг
	ns := nats.NewNats(&cfg.Nats)
	fmt.Println("Nats server is running, successfully connected")

	//подключаем бд
	postgresConnect, err := postgres.ConnectionPG(&cfg.PG)

	if err != nil {
		fmt.Printf("Error connecting to Postgresql: %v", err)
	}
	defer postgresConnect.Close()

	// соед-е бд
	repo := repository.NewRepository(postgresConnect)

	dbCreatErr := repo.CreateTable()

	if dbCreatErr != nil {
		fmt.Printf("Error creating table: %v", dbCreatErr)
	}

	// созд-е кэша (из бд)
	orderCache := cache.NewCache(repo)
	orderCache.Preload()

	//публикация(отправка заказов)в натс каждые 30 сек
	go func() {
		for {
			order := generator.GenerateOrder()
			fmt.Println("Order sent") //заказ отправлен
			err := ns.Publish(order)  //

			if err != nil {
				fmt.Printf("Error while publishing: %v\n", err)
			}

			time.Sleep(30 * time.Second)
		}
	}()

	//подписка(получение сообщений от натса и сохраняет их в кэш)
	go func() {
		for {
			mes, err := ns.Subscribe()
			fmt.Println("Order received") //заказ получен
			if err != nil {
				fmt.Printf("Error while subscribing: %v", err)
			}

			if err != nil {
				fmt.Printf("Error while Unmarshaling: %v", err)
			}

			orderCache.CreateCache(*mes)

			time.Sleep(30 * time.Second)
		}
	}()

	httpServer := httpserver.NewServer()
	orderController := controller.NewOrderController(orderCache)

	serverStartingError := httpServer.Start(orderController.GetOrderController, orderController.GetAllOrders)

	if serverStartingError != nil {
		log.Fatalf("Error at server starting: %v", serverStartingError)
	}
}
