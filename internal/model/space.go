package model

import (
	"gorm.io/gorm"
	"time"
	"yema.dev/app/model/field"
)

type Space struct {
	ID     int64        `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
	UserId int64        `gorm:"column:user_id;notNull;comment:所属用户id" json:"user_id"`
	Name   string       `gorm:"column:name;size:100;notNull;comment:空间名" json:"name"`
	Status field.Status `gorm:"column:status;notNull;default:0;comment:状态" json:"status"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;notNull" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;notNull" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	User     User       ` json:"user"`
	Projects []*Project `json:"projects"`
}
