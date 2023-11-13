package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"yema.dev/api/space"
	"yema.dev/internal/errcode"
	"yema.dev/internal/model"
)

type SpaceRepository interface {
	List(ctx context.Context, req *space.ListReq) (total int64, res []*model.Space, err error)
	Create(ctx context.Context, m *model.Space) error
	Update(ctx context.Context, user *model.Space) error
	UpdateFields(ctx context.Context, user *model.Space, fields ...string) error
	GetByID(ctx context.Context, id int64) (*model.Space, error)
	DeleteByID(ctx context.Context, id int64) error
	GetProjectCount(ctx context.Context, id int64) int64
}

func NewSpaceRepository(r *Repository) SpaceRepository {
	return &spaceRepository{
		Repository: r,
	}
}

type spaceRepository struct {
	*Repository
}

func (r *spaceRepository) List(ctx context.Context, req *space.ListReq) (total int64, list []*model.Space, err error) {
	err = r.DB(ctx).Model(&model.Space{}).Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = r.DB(ctx).Model(&model.Space{}).Scopes(req.PageQuery()).Preload("User").Find(&list).Error
	return
}

func (r *spaceRepository) Create(ctx context.Context, m *model.Space) error {
	if err := r.DB(ctx).Create(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *spaceRepository) Update(ctx context.Context, m *model.Space) error {
	if err := r.DB(ctx).Save(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *spaceRepository) UpdateFields(ctx context.Context, m *model.Space, fields ...string) error {
	if err := r.DB(ctx).Select(fields).Where("id = ?", m.ID).Updates(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *spaceRepository) GetByID(ctx context.Context, id int64) (*model.Space, error) {
	var m model.Space
	if err := r.DB(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (r *spaceRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.DB(ctx).Delete(&model.Space{}, id).Error
}

func (r *spaceRepository) GetProjectCount(ctx context.Context, id int64) int64 {
	return r.DB(ctx).Model(&model.Space{ID: id}).Association("Projects").Count()
}
