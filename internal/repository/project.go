package repository

import (
	"context"
	"yema.dev/api/project"
	"yema.dev/internal/model"
)

type ProjectRepository interface {
	List(ctx context.Context, req *project.ListReq) (total int64, list []*model.Project, err error)
	Create(ctx context.Context, m *model.Project) error
	Update(ctx context.Context, user *model.Project) error
	UpdateFields(ctx context.Context, m *model.Project, fields ...string) error
	GetByID(ctx context.Context, id int64) (*model.Project, error)
	DeleteByID(ctx context.Context, id int64) error
	ClearServers(ctx context.Context, id int64) error
}

func NewProjectRepository(r *Repository) ProjectRepository {
	return &projectRepository{
		Repository: r,
	}
}

type projectRepository struct {
	*Repository
}

func (r *projectRepository) List(ctx context.Context, req *project.ListReq) (total int64, list []*model.Project, err error) {
	where := model.Project{SpaceId: req.SpaceId}
	if req.EnvironmentId > 0 {
		where.EnvironmentId = req.EnvironmentId
	}
	_db := r.DB(ctx).Model(&where).Where(where)
	err = _db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = _db.Scopes(req.PageQuery()).
		Preload("Space").
		Preload("Environment").Find(&list).Error
	return
}

func (r *projectRepository) Create(ctx context.Context, m *model.Project) error {
	if err := r.DB(ctx).Create(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *projectRepository) Update(ctx context.Context, m *model.Project) error {
	if err := r.DB(ctx).Save(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *projectRepository) UpdateFields(ctx context.Context, m *model.Project, fields ...string) error {
	if err := r.DB(ctx).Select(fields).Where("id = ?", m.ID).
		Updates(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *projectRepository) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	var m model.Project
	if err := r.DB(ctx).Where("id = ?", id).Preload("Servers").Preload("Environment").First(&m).Error; err != nil {
		//if errors.Is(err, gorm.ErrRecordNotFound) {
		//	return nil, errcode.ErrNotFound
		//}
		return nil, err
	}
	return &m, nil
}

func (r *projectRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.DB(ctx).Delete(&model.Project{}, id).Error
}

// ClearServers 清除所有跟server 的绑定关系
func (r *projectRepository) ClearServers(ctx context.Context, id int64) error {
	return r.DB(ctx).Model(&model.Project{ID: id}).Association("Servers").Clear()
}
