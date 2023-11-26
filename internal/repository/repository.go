package repository

import (
	"context"
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
	CREATE TABLE IF NOT EXISTS "order" (
		order_uid VARCHAR(50) PRIMARY KEY ,
		track_number VARCHAR(50) NOT NULL UNIQUE,
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
        INSERT INTO "order" (order_uid, track_number, entry, delivery_info, payment_info, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    `, order.OrderUid, order.TrackNumber, order.Entry, order.Delivery, order.Payment,
		order.Items, order.Locale, order.InternalSignature, order.CustomerId,
		order.DeliveryService, order.ShardKey, order.SmId, order.DateCreated, order.OofShard)

	return err
}
