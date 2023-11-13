package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
	"yema.dev/internal/model/field"
)

type Server struct {
	ID          int64        `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
	SpaceId     int64        `gorm:"column:space_id;index;notNull;comment:所属空间" json:"space_id"`
	Name        string       `gorm:"column:name;type:string;size:100;notNull;comment:名称" json:"name"`
	User        string       `gorm:"column:user;type:string;size:100;notNull;comment:用户名" json:"user"`
	Host        string       `gorm:"column:host;type:string;size:100;notNull;comment:主机" json:"host"`
	Port        int          `gorm:"column:port;notNull;default:22;comment:端口" json:"port"`
	Status      field.Status `gorm:"column:status;notNull;default:0;comment:状态" json:"status"`
	Description string       `gorm:"column:description;type:string;size:500;notNull;default:'';comment:简介说明" json:"description"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;notNull" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;notNull" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	Projects []Project `gorm:"many2many:project_server" json:"projects,omitempty'"`
	Tasks    []Task    `gorm:"many2many:task_server" json:"tasks,omitempty"`
}

func (receiver *Server) Hostname() string {
	return fmt.Sprintf("%s@%s", receiver.User, receiver.Host)
}
