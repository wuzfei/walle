package server

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"yema.dev/internal/model"
)

type Migrate struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewMigrate(db *gorm.DB, log *zap.Logger) *Migrate {
	return &Migrate{
		db:  db,
		log: log,
	}
}
func (m *Migrate) Start(ctx context.Context) error {
	if err := m.db.AutoMigrate(&model.User{}); err != nil {
		m.log.Error("user migrate error", zap.Error(err))
		return err
	}
	m.log.Info("AutoMigrate success")
	os.Exit(0)
	return nil
}
func (m *Migrate) Stop(ctx context.Context) error {
	m.log.Info("AutoMigrate stop")
	return nil
}
