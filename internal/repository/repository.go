package repository

import (
	"context"
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
		order_uid VARCHAR(50) PRIMARY KEY NOT NULL,
		track_number VARCHAR(50) NOT NULL,
		entry VARCHAR(55) NOT NULL,
		delivery_info JSONB,
		payment_info JSONB,
		items JSONB,
		locale VARCHAR(2) NOT NULL,
		internal_signature VARCHAR(50) NOT NULL,
		customer_id VARCHAR(50) NOT NULL,
		delivery_service VARCHAR(50) NOT NULL,
		shardkey VARCHAR(50) NOT NULL,
		sm_id BIGINT CHECK (sm_id > 0),
		date_created TIMESTAMP NOT NULL,
		oof_shard VARCHAR(50) NOT NULL
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
		err := res.Scan(
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
			fmt.Printf("Error at parsing preloading: %v\n", err)
		}

		orders = append(orders, execOrder)
	}

	if err := res.Err(); err != nil {
		fmt.Printf("Error at final getting orders: %v\n", err)
		return nil, err
	}
	fmt.Println("Orders preloaded DB -> Cache")
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
