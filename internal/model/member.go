package model

import (
	"time"
)

type Member struct {
	ID      int64  `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
	UserId  int64  `gorm:"column:user_id;notNull;uniqueIndex:space_user;comment:用户" json:"user_id"`
	SpaceId int64  `gorm:"column:space_id;notNull;uniqueIndex:space_user;comment:空间" json:"space_id"`
	Role    string `gorm:"column:role;size:20;notNull;comment:角色" json:"role"`

	CreatedAt time.Time `gorm:"column:created_at;type:datetime;notNull" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;notNull" json:"updated_at"`

	Space Space
	User  User
}

type Role string

const RoleSuper Role = "super"
const RoleOwner Role = "owner"
const RoleMaster Role = "master"
const RoleDeveloper Role = "developer"
const RoleVisitor Role = "visitor"

var roleLevel = map[Role]int{
	RoleVisitor:   1,
	RoleDeveloper: 1 << 2,
	RoleMaster:    1 << 3,
	RoleOwner:     1 << 4,
	RoleSuper:     1 << 5,
}

func (r Role) Level() int {
	if v, ok := roleLevel[r]; ok {
		return v
	}
	return 0
}
