package repository

import (
	"context"
	"yema.dev/api/deploy"
	"yema.dev/internal/model"
)

type DeployRepository interface {
	List(ctx context.Context, req *deploy.ListReq) (total int64, list []*model.Task, err error)
	Create(ctx context.Context, m *model.Task) error
	Update(ctx context.Context, user *model.Task) error
	UpdateFields(ctx context.Context, m *model.Task, fields ...string) error
	GetByID(ctx context.Context, id int64) (*model.Task, error)
	DeleteByID(ctx context.Context, id int64) error
}

func NewDeployRepository(r *Repository) DeployRepository {
	return &deployRepository{
		Repository: r,
	}
}

type deployRepository struct {
	*Repository
}

func (r *deployRepository) List(ctx context.Context, req *deploy.ListReq) (total int64, list []*model.Task, err error) {
	_db := r.DB(ctx).Model(&model.Task{}).Where("space_id=?", req.SpaceId)
	err = _db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = _db.Scopes(req.PageQuery()).
		Preload("User").
		Preload("Project").
		Preload("Environment").
		Order("id desc").
		Find(&list).Error
	return
}

func (r *deployRepository) Create(ctx context.Context, m *model.Task) error {
	if err := r.DB(ctx).Create(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *deployRepository) Update(ctx context.Context, m *model.Task) error {
	if err := r.DB(ctx).Save(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *deployRepository) UpdateFields(ctx context.Context, m *model.Task, fields ...string) error {
	if err := r.DB(ctx).Select(fields).Where("id = ?", m.ID).
		Updates(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *deployRepository) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	var m model.Task
	if err := r.DB(ctx).Where("id = ?", id).Preload("Servers").First(&m).Error; err != nil {
		//if errors.Is(err, gorm.ErrRecordNotFound) {
		//	return nil, errcode.ErrNotFound
		//}
		return nil, err
	}
	return &m, nil
}

func (r *deployRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.DB(ctx).Delete(&model.Task{}, id).Error
}
