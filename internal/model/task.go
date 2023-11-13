package model

import (
	"database/sql"
	"time"
)

const (
	TaskStatusWaiting         = 1 //新建提交，等待审核
	TaskStatusAudit           = 2 //审核通过
	TaskStatusReject          = 3 //审核拒绝
	TaskStatusRelease         = 4 //上线发布中
	TaskStatusReleaseFail     = 5 //上线失败
	TaskStatusReleasePartFail = 6 //部分服务器失败
	TaskStatusFinish          = 7 //上线完成
)

type Task struct {
	ID            int64 `gorm:"column:id" json:"id"`
	SpaceId       int64 `gorm:"column:space_id;index;notNull;comment:所属空间" json:"space_id"`
	ProjectId     int64 `gorm:"column:project_id;notNull;comment:所属项目" json:"project_id"`
	UserId        int64 `gorm:"column:user_id;notNull;comment:所属用户" json:"user_id"`
	EnvironmentId int64 `gorm:"column:environment_id;notNull;comment:所属环境" json:"environment_id"`

	Name        string       `gorm:"column:Name;size:100;notNull;comment:名称" json:"name"`
	Status      int8         `gorm:"column:status;notNull;default:0;comment:状态" json:"status"`
	Version     string       `gorm:"column:version;type:string;size:100;notNull;default:'';comment:版本号" json:"version"`
	PrevVersion string       `gorm:"column:prev_version;type:string;size:100;notNull;default:'';comment:上一个版本号" json:"prev_version"`
	CommitId    string       `gorm:"column:commit_id;type:string;size:100;notNull;default:'';comment:commit哈希" json:"commit_id"`
	Branch      string       `gorm:"column:branch;type:string;size:100;notNull;default:'';comment:分支" json:"branch"`
	Tag         string       `gorm:"column:tag;type:string;size:100;notNull;default:'';comment:tag" json:"tag"`
	IsRollback  int8         `gorm:"column:is_rollback;notNull;default:0;comment:是否回滚" json:"is_rollback"`
	LastError   string       `gorm:"column:last_error;type:string;notNull;default:'';comment:最后错误" json:"last_error"`
	AuditUserId int64        `gorm:"column:audit_user_id;notNull;default:0;审核员" json:"audit_user_id"`
	AuditTime   sql.NullTime `gorm:"column:audit_time;type:datetime;最后审核操作时间" json:"audit_time"`

	CreatedAt time.Time `gorm:"column:created_at;type:datetime;notNull" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;notNull" json:"updated_at"`

	Project     Project     `json:"project,omitempty"`
	User        User        `json:"user,omitempty"`
	Space       Space       `json:"space,omitempty"`
	Environment Environment `json:"environment,omitempty"`
	Servers     []Server    `gorm:"many2many:task_server" json:"servers"`
}
