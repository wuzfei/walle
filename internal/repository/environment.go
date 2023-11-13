package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"yema.dev/api/environment"
	"yema.dev/internal/errcode"
	"yema.dev/internal/model"
)

type EnvironmentRepository interface {
	List(ctx context.Context, req *environment.ListReq) (total int64, res []*model.Environment, err error)
	Create(ctx context.Context, m *model.Environment) error
	Update(ctx context.Context, user *model.Environment) error
	UpdateFields(ctx context.Context, m *model.Environment, fields ...string) error
	GetByID(ctx context.Context, id int64) (*model.Environment, error)
	DeleteByID(ctx context.Context, id int64) error
	GetProjectCount(ctx context.Context, id int64) int64
}

func NewEnvironmentRepository(r *Repository) EnvironmentRepository {
	return &environmentRepository{
		Repository: r,
	}
}

type environmentRepository struct {
	*Repository
}

func (r *environmentRepository) List(ctx context.Context, req *environment.ListReq) (total int64, list []*model.Environment, err error) {
	_db := r.DB(ctx).Model(&model.Environment{}).Where("space_id = ?", req.SpaceId)
	if req.Kw != "" {
		_db = _db.Where("name like ", "%"+req.Kw+"%")
	}
	err = _db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = _db.Scopes(req.PageQuery()).Preload("Space").Find(&list).Error
	return
}

func (r *environmentRepository) Create(ctx context.Context, m *model.Environment) error {
	if err := r.DB(ctx).Create(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *environmentRepository) Update(ctx context.Context, m *model.Environment) error {
	if err := r.DB(ctx).Save(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *environmentRepository) UpdateFields(ctx context.Context, m *model.Environment, fields ...string) error {
	if err := r.DB(ctx).Select(fields).Where("id = ?", m.ID).
		Updates(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *environmentRepository) GetByID(ctx context.Context, id int64) (*model.Environment, error) {
	var m model.Environment
	if err := r.DB(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (r *environmentRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.DB(ctx).Delete(&model.Environment{}, id).Error
}

func (r *environmentRepository) GetProjectCount(ctx context.Context, id int64) int64 {
	return r.DB(ctx).Model(&model.Environment{ID: id}).Association("Projects").Count()
}
