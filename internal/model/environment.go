package model

import (
	"gorm.io/gorm"
	"time"
	"yema.dev/app/model/field"
)

type Environment struct {
	ID          int64        `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
	SpaceId     int64        `gorm:"column:space_id;notNull;comment:所属空间" json:"space_id"`
	Name        string       `gorm:"column:name;type:string;size:100;notNull;comment:名称" json:"name"`
	Status      field.Status `gorm:"column:status;notNull;default:0;comment:状态" json:"status"`
	Description string       `gorm:"column:description;type:string;size:500;notNull;default:'';comment:简介说明" json:"description"`
	Color       string       `gorm:"column:color;size:10;notNull;default:'';comment:主题色" json:"color"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;notNull" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;notNull" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	Space    Space      `json:"space"`
	Projects []*Project `json:"projects"`
}
