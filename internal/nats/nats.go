package nats

import (
	"encoding/json"
	"fmt"
	"test_wb/config"
	"test_wb/internal/models"
	"time"

	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
)

type Nats struct {
	config *config.Nats
	sc     stan.Conn
	nc     *nats.Conn
}

func NewNats(ncfg *config.Nats) *Nats {
	natsUrl := fmt.Sprintf("nats://%s:%s", ncfg.Host, ncfg.Port)

	// Подключение к серверу NATS
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		fmt.Printf("Can't connect to Nats: %v", err)
		return nil
	}
	defer nc.Close()

	// Подключение к кластеру NATS Streaming
	sc, err := stan.Connect(ncfg.Cluster, ncfg.Client) // stan.NatsConn(nc) - это метод или функция, используемая для создания подключения к NATS Streaming через уже установленное соединение NATS.
	if err != nil {
		fmt.Printf("Can't connect to Nats-Streaming: %v", err)
		return nil
	}
	defer sc.Close()

	return &Nats{ncfg, sc, nc}
}

// отправка сообщений (публикация сообщ о заказе)
func (ns *Nats) Publish(message *models.Order) error {

	ord, err := json.MarshalIndent(message, "", "\t")

	if err != nil {
		fmt.Printf("Error at marshaling new order: %v", err)
	}

	return ns.sc.Publish(ns.config.Topic, ord) //публикует в натс (публиш тут функция с библиотеки, а не метод)
}

// получение сообщений (подписка на сообщения)
func (ns *Nats) Subscribe() (*models.Order, error) {

	var rc models.Order

	ch := make(chan *models.Order)
	/*Подписывается на тему в NATS через ns.sc.Subscribe, который
	ожидает получения сообщения)*/
	_, err := ns.sc.Subscribe(ns.config.Topic, func(mes *stan.Msg) {

		// Когда сообщ поступило - распаковывается с помощью json.Unmarshal в объект models.Order
		err := json.Unmarshal(mes.Data, &rc)

		if err != nil {
			fmt.Printf("Error at Unmarshaling: %v", err)
		}

		ch <- &rc //cообщ отправляется в канал
	})

	if err != nil {
		fmt.Printf("Error at subscription: %v", err)
	}

	/*ожиданиe сообщения от канала.
	блокируется до тех пор, пока не получит сообщение, и возвращает распакованное
	сообщение о заказе или же возвращает ошибку stan.ErrTimeout,
	если происходит таймаут (60 секунд в данном случае).*/
	select {
	case rc := <-ch:
		return rc, nil
	case <-time.After(60 * time.Second):
		return nil, stan.ErrTimeout //ОШИБКА:В КАНАЛЕ ПУСТО
	}
}
