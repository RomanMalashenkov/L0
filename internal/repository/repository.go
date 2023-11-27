package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"test_wb/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repo {
	return &Repo{
		pool: pool,
	}
}

func (r *Repo) CreateTable() error {
	_, err := r.pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS orders (
		order_uid VARCHAR(50) PRIMARY KEY,
		track_number VARCHAR(50),
		entry VARCHAR(55),
		delivery_info JSONB,
		payment_info JSONB,
		items JSONB,
		locale VARCHAR(2),
		internal_signature VARCHAR(50),
		customer_id VARCHAR(50),
		delivery_service VARCHAR(50),
		shardkey VARCHAR(50),
		sm_id VARCHAR(50),
		date_created VARCHAR(50),
		oof_shard VARCHAR(50)
	)
`)
	return err
}

func (r *Repo) SaveOrder(order models.Order) error {
	_, err := r.pool.Exec(context.Background(), `
        INSERT INTO orders (order_uid, track_number, entry, delivery_info, payment_info, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    `, order.OrderUid, order.TrackNumber, order.Entry, order.Delivery, order.Payment,
		order.Items, order.Locale, order.InternalSignature, order.CustomerId,
		order.DeliveryService, order.ShardKey, order.SmId, order.DateCreated, order.OofShard)

	return err
}

func (r *Repo) GetALl() ([]models.Order, error) {
	fmt.Println("Attempting to retrieve orders from the database...") //sdf

	res, err := r.pool.Query(context.Background(),
		`SELECT * FROM orders`)

	if err != nil {
		fmt.Printf("Error at getting all orders: %v\n", err)
		return nil, err
	}

	defer res.Close()

	var orders []models.Order
	for res.Next() {
		var execOrder models.Order
		var deliveryJSON, paymentJSON, itemsJSON []byte // JSON в виде []byte

		// Сканирование всех столбцов, включая JSONB
		err := res.Scan(
			&execOrder.OrderUid,
			&execOrder.TrackNumber,
			&execOrder.Entry,
			&deliveryJSON,
			&paymentJSON,
			&itemsJSON,
			&execOrder.Locale,
			&execOrder.InternalSignature,
			&execOrder.CustomerId,
			&execOrder.DeliveryService,
			&execOrder.ShardKey,
			&execOrder.SmId,
			&execOrder.DateCreated,
			&execOrder.OofShard,
		)

		if err != nil {
			fmt.Printf("Error at parsing preloading: %v\n", err)
			continue // Продолжить цикл даже при ошибке чтения строки
		}

		/////////////////////////
		// Распаковка JSON в соответствующие структуры
		var delivery models.Delivery
		err = json.Unmarshal([]byte(deliveryJSON), &delivery)
		if err != nil {
			fmt.Printf("Error unmarshaling delivery_info: %v\n", err)
			continue
		}
		execOrder.Delivery = delivery

		var payment models.Payment
		err = json.Unmarshal([]byte(paymentJSON), &payment)
		if err != nil {
			fmt.Printf("Error unmarshaling payment_info: %v\n", err)
			continue
		}
		execOrder.Payment = payment

		var items []models.Item
		err = json.Unmarshal([]byte(itemsJSON), &items)
		if err != nil {
			fmt.Printf("Error unmarshaling items: %v\n", err)
			continue
		}
		execOrder.Items = items

		///////////////////////////

		orders = append(orders, execOrder)
	}

	if err := res.Err(); err != nil {
		fmt.Printf("Error at final getting orders: %v\n", err)
		return nil, err
	}
	//fmt.Println("Orders preloaded DB -> Cache")
	fmt.Printf("Retrieved %d orders from the database\n", len(orders)) // Добавлено сообщение о количестве полученных заказов

	return orders, nil
}

func (r *Repo) GetOrder(uid string) (models.Order, error) {

	var execOrder models.Order

	err := r.pool.QueryRow(context.Background(), `SELECT * FROM orders WHERE order_uid = $1`, uid).
		Scan(
			&execOrder.OrderUid,
			&execOrder.TrackNumber,
			&execOrder.Entry,
			&execOrder.Delivery,
			&execOrder.Payment,
			&execOrder.Items,
			&execOrder.Locale,
			&execOrder.InternalSignature,
			&execOrder.CustomerId,
			&execOrder.DeliveryService,
			&execOrder.ShardKey,
			&execOrder.SmId,
			&execOrder.DateCreated,
			&execOrder.OofShard,
		)

	if err != nil {
		return execOrder, err
	}

	return execOrder, nil

}
