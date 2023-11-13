package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"yema.dev/api/server"
	"yema.dev/internal/model"
)

type ServerRepository interface {
	List(ctx context.Context, req *server.ListReq) (total int64, res []*model.Server, err error)
	Create(ctx context.Context, m *model.Server) error
	Update(ctx context.Context, user *model.Server) error
	UpdateFields(ctx context.Context, m *model.Server, fields ...string) error
	GetByID(ctx context.Context, id int64) (*model.Server, error)
	DeleteByID(ctx context.Context, id int64) error
	FindByHostIp(ctx context.Context, spaceId int64, user, host string, port int) (m model.Server, err error)
	ClearProjects(ctx context.Context, id int64) error
}

func NewServerRepository(r *Repository) ServerRepository {
	return &serverRepository{
		Repository: r,
	}
}

type serverRepository struct {
	*Repository
}

func (r *serverRepository) List(ctx context.Context, req *server.ListReq) (total int64, list []*model.Server, err error) {
	_db := r.DB(ctx).Model(&model.Server{}).Where("space_id = ? ", req.SpaceId)
	err = _db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = _db.Scopes(req.PageQuery()).Find(&list).Error
	return
}

func (r *serverRepository) Create(ctx context.Context, m *model.Server) error {
	if err := r.DB(ctx).Create(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *serverRepository) Update(ctx context.Context, m *model.Server) error {
	if err := r.DB(ctx).Save(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *serverRepository) UpdateFields(ctx context.Context, m *model.Server, fields ...string) error {
	if err := r.DB(ctx).Select(fields).Where("id = ?", m.ID).
		Updates(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *serverRepository) GetByID(ctx context.Context, id int64) (*model.Server, error) {
	var m model.Server
	if err := r.DB(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		//if errors.Is(err, gorm.ErrRecordNotFound) {
		//	return nil, errcode.ErrNotFound
		//}
		return nil, err
	}
	return &m, nil
}

func (r *serverRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.DB(ctx).Delete(&model.Server{}, id).Error
}

// FindByHostIp aa
func (r *serverRepository) FindByHostIp(ctx context.Context, spaceId int64, user, host string, port int) (m model.Server, err error) {
	err = r.DB(ctx).Where(&model.Server{SpaceId: spaceId, User: user, Host: host, Port: port}).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// ClearProjects 清除所有跟project 的绑定关系
func (r *serverRepository) ClearProjects(ctx context.Context, id int64) error {
	return r.DB(ctx).Model(&model.Server{ID: id}).Association("Projects").Clear()
}
