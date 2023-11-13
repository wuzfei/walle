package model

import (
	"gorm.io/gorm"
	"time"
	"yema.dev/internal/model/field"
)

const SuperUserId = 1

func IsSuperUser(userId int64) bool {
	return userId == SuperUserId
}

type User struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`

	Username      string       `gorm:"column:username;type:string;size:100;notNull;default:'';comment:用户名" json:"username"`
	Email         string       `gorm:"column:email;uniqueIndex;type:string;size:100;notNull;comment:邮箱" json:"email"`
	Password      []byte       `gorm:"column:password;type:string;size:200;notNull;comment:密码" json:"-"`
	Status        field.Status `gorm:"column:status;size:1;notNull;default:0;comment:状态" json:"status"`
	RememberToken string       `gorm:"remember_token;type:string;size:500;notNull;default:'';comment:记住密码token" json:"-"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;notNull" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;notNull" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	Spaces []*Space `gorm:"foreignKey:user_id"`
}

func (u *User) IsSuperUser() bool {
	return IsSuperUser(u.ID)
}
