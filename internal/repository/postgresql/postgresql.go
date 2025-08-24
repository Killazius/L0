package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/Killazius/L0/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Close() {
	r.DB.Close()
}

func (r *Repository) Create(ctx context.Context, order domain.Order) error {

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	_, err = tx.Exec(ctx, `
        INSERT INTO "orders" (
            order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO deliveries (
            order_uid, name, phone, zip, city, address, region, email
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO payments (
            order_uid, transaction, request_id, currency, provider,
            amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	for _, item := range order.Items {
		_, err = tx.Exec(ctx, `
            INSERT INTO items (
                order_uid, chrt_id, track_number, price, rid, name,
                sale, size, total_price, nm_id, brand, status
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        `,
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to insert item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, orderUID string) (*domain.Order, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	order, err := r.getOrder(ctx, tx, orderUID)
	if err != nil {
		return nil, err
	}

	delivery, err := r.getDelivery(ctx, tx, orderUID)
	if err != nil {
		return nil, err
	}
	order.Delivery = *delivery

	payment, err := r.getPayment(ctx, tx, orderUID)
	if err != nil {
		return nil, err
	}
	order.Payment = *payment

	items, err := r.getItems(ctx, tx, orderUID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

func (r *Repository) getOrder(ctx context.Context, tx pgx.Tx, orderUID string) (*domain.Order, error) {
	var order domain.Order

	query := `
		SELECT 
			order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders 
		WHERE order_uid = $1
	`

	err := tx.QueryRow(ctx, query, orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("order with UID %s not found", orderUID)
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

func (r *Repository) getDelivery(ctx context.Context, tx pgx.Tx, orderUID string) (*domain.Delivery, error) {
	var delivery domain.Delivery

	query := `
		SELECT 
			name, phone, zip, city, address, region, email
		FROM deliveries 
		WHERE order_uid = $1
	`

	err := tx.QueryRow(ctx, query, orderUID).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("delivery for order UID %s not found", orderUID)
		}
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}

	return &delivery, nil
}

func (r *Repository) getPayment(ctx context.Context, tx pgx.Tx, orderUID string) (*domain.Payment, error) {
	var payment domain.Payment

	query := `
		SELECT 
			transaction, request_id, currency, provider, amount,
			payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payments 
		WHERE order_uid = $1
	`

	err := tx.QueryRow(ctx, query, orderUID).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("payment for order UID %s not found", orderUID)
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &payment, nil
}

func (r *Repository) getItems(ctx context.Context, tx pgx.Tx, orderUID string) ([]domain.Item, error) {
	var items []domain.Item

	query := `
		SELECT 
			chrt_id, track_number, price, rid, name, sale,
			size, total_price, nm_id, brand, status
		FROM items 
		WHERE order_uid = $1
		ORDER BY id
	`

	rows, err := tx.Query(ctx, query, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating items: %w", err)
	}

	return items, nil
}
