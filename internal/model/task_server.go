package model

const (
	TaskServerStatus        = 0
	TaskServerStatusSuccess = 1
	TaskServerStatusFail    = 2
)

type TaskServer struct {
	TaskId   int64  `gorm:"column:task_id;primaryKey:taskServerIdx;notNull" json:"task_id"`
	ServerId int64  `gorm:"column:server_id;primaryKey:taskServerIdx;notNull" json:"server_id"`
	Status   int8   `gorm:"column:status;notNull;default:0" json:"status"`
	Err      string `gorm:"column:err;size:1000;notNull;default:''" json:"err"`
}

func (t *TaskServer) TableName() string {
	return "task_server"
}
