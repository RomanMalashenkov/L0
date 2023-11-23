package main

import (
	"encoding/json"
	"fmt"
	"time"

	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
)

type NNats struct {
	config *Nats
	sc     stan.Conn
	nc     *nats.Conn
}

func NewNats(ncfg *Nats) *NNats {
	natsUrl := fmt.Sprintf("nats://%s:%s", ncfg.Host, ncfg.Port)

	// Подключение к серверу NATS
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		fmt.Printf("Can't connect to Nats: %v", err)
		return nil
	}
	defer nc.Close()

	// Подключение к кластеру NATS Streaming
	sc, err := stan.Connect(ncfg.Cluster, ncfg.Client, stan.NatsConn(nc)) // stan.NatsConn(nc) - это метод или функция, используемая для создания подключения к NATS Streaming через уже установленное соединение NATS.
	if err != nil {
		fmt.Printf("Can't connect to Nats-Streaming: %v", err)
		return nil
	}
	defer sc.Close()

	return &NNats{ncfg, sc, nc}
}

func (ns *Nats) Publish(message *Order) error {

	ord, err := json.MarshalIndent(message, "", "\t")

	if err != nil {
		fmt.Printf("Error at marshaling new order: %v", err)
	}

	return ns.sc.Publish(ns.Topic, ord)
}

// два метода для структуры
func (ns *Nats) Subscribe() (*Order, error) {

	var rc models.Order

	ch := make(chan *models.Order)

	_, err := ns.sc.Subscribe(ns.config.Topic, func(mes *stan.Msg) {

		err := json.Unmarshal(mes.Data, &rc)

		if err != nil {
			fmt.Printf("Error at Unmarshaling: %v", err)
		}

		ch <- &rc
	})

	if err != nil {
		fmt.Printf("Error at subscription: %v", err)
	}

	select {
	case rc := <-ch:
		return rc, nil
	case <-time.After(60 * time.Second):
		return nil, stan.ErrTimeout
	}

}
