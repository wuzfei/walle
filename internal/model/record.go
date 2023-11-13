package model

import (
	"errors"
	"time"
	"yema.dev/app/model/field"
)

const (
	RecordTypeDefault = iota
	RecordTypePrevDeploy
	RecordTypeDeploy
	RecordTypePostDeploy
	RecordTypePrevRelease
	RecordTypeRelease
	RecordTypePostRelease
)

type Record struct {
	ID       int64                `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
	UserId   int64                `gorm:"column:user_id;notNull;comment:操作用户" json:"user_id"`
	ServerId int64                `gorm:"column:server_id;notNull;default:0;comment:服务器id" json:"server_id"`
	TaskId   int64                `gorm:"column:task_id;notNull;default:0;comment:所属任务" json:"task_id"`
	Envs     field.Slices[string] `gorm:"column:envs;notNull;default:'';comment:执行环境变量" json:"envs"`
	RunTime  int64                `gorm:"column:run_time;notNull;default:0;comment:执行时间(ms)" json:"run_time"`
	Status   int                  `gorm:"column:status;notNull;default:0;comment:执行状态,linux退出码" json:"status"`
	Command  string               `gorm:"column:command;notNull;default:'';comment:执行命令" json:"command"`
	Output   string               `gorm:"column:output;notNull;default:'';comment:输出结果" json:"output"`
	Step     int8                 `gorm:"column:step;notNull;default:0;comment:流程步骤" json:"step"`

	Server Server `json:"server"`

	CreatedAt time.Time `gorm:"column:created_at;type:datetime;notNull" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;notNull" json:"updated_at"`
}

func (receiver *Record) ErrExit() bool {
	return receiver.Status != 0
}

func (receiver *Record) ExecError() error {
	if receiver.Status != 0 {
		return errors.New(receiver.Output)
	}
	return nil
}
