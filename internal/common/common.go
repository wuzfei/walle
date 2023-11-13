package common

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SpaceWithId struct {
	SpaceId, ID int64
}

type Service struct {
	log *zap.Logger
	db  *gorm.DB
}

func NewService(log *zap.Logger, db *gorm.DB) *Service {
	onceService.Do(func() {
		service = &Service{log: log, db: db}
	})
	return service
}

// Statistics 发布数据统计
func (srv *Service) Statistics() StatisticsRes {
	return StatisticsRes{}
}

// ServerInfo 系统系统
func (srv *Service) ServerInfo() (*ServerInfo, error) {
	return getServerInfo()
}

// WaitAudit 待审核列表
func (srv *Service) WaitAudit() {

}

// Release 最近发布成功消息
func (srv *Service) Release() {

}
