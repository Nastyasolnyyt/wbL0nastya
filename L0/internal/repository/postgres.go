package repository

import (
	"context"
	"fmt"
	"myapp/internal/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type OrderRepository struct {
	db *sqlx.DB
}

func New(dbURL string) (*OrderRepository, error) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	return &OrderRepository{db: db}, nil
}

func (r *OrderRepository) Close() {
	r.db.Close()
}

func (r *OrderRepository) SaveOrder(ctx context.Context, order model.Order) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExecContext(ctx, `
        INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
        VALUES (:order_uid, :name, :phone, :zip, :city, :address, :region, :email)`,
		order.Delivery)
	if err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	_, err = tx.NamedExecContext(ctx, `
        INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES (:transaction, :request_id, :currency, :provider, :amount, :payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee)`,
		order.Payment)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	_, err = tx.NamedExecContext(ctx, `
        INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES (:order_uid, :track_number, :entry, :locale, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard)`,
		order)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	for _, item := range order.Items {
		_, err = tx.NamedExecContext(ctx, `
            INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES (:order_uid, :chrt_id, :track_number, :price, :rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status)`,
			item)
		if err != nil {
			return fmt.Errorf("failed to insert item: %w", err)
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order

	err := r.db.SelectContext(ctx, &orders, "SELECT * FROM orders")
	if err != nil {
		return nil, err
	}

	for i := range orders {

		err = r.db.GetContext(ctx, &orders[i].Delivery, "SELECT * FROM deliveries WHERE order_uid=$1", orders[i].OrderUID)
		if err != nil {
			return nil, err
		}

		err = r.db.GetContext(ctx, &orders[i].Payment, "SELECT * FROM payments WHERE transaction=$1", orders[i].OrderUID)
		if err != nil {
			return nil, err
		}

		err = r.db.SelectContext(ctx, &orders[i].Items, "SELECT * FROM items WHERE order_uid=$1", orders[i].OrderUID)
		if err != nil {
			return nil, err
		}
	}

	return orders, nil
}
