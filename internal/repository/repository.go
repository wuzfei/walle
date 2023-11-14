package repository

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const ctxTxKey = "TxKey"

type Repository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewRepository(db *gorm.DB, log *zap.Logger) *Repository {
	return &Repository{
		db:  db,
		log: log.Named("repository"),
	}
}

type Transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewTransaction(r *Repository) Transaction {
	return r
}

// DB return tx
// If you need to create a Transaction, you must call DB(ctx) and Transaction(ctx,fn)
func (r *Repository) DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ctxTxKey)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return r.db.WithContext(ctx)
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, ctxTxKey, tx)
		return fn(ctx)
	})
}
