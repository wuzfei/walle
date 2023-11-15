package repository

import (
	"context"
	"gorm.io/gorm/clause"
	"yema.dev/api/member"
	"yema.dev/internal/model"
)

type MemberRepository interface {
	List(ctx context.Context, req *member.ListReq) (total int64, res []*model.Member, err error)
	Store(ctx context.Context, m *model.Member) error
	GetByID(ctx context.Context, id int64) (*model.Member, error)
	DeleteByID(ctx context.Context, id int64) error
}

func NewMemberRepository(r *Repository) MemberRepository {
	return &memberRepository{
		Repository: r,
	}
}

type memberRepository struct {
	*Repository
}

func (r *memberRepository) List(ctx context.Context, req *member.ListReq) (total int64, list []*model.Member, err error) {
	_db := r.DB(ctx).Model(&model.Member{}).Where("space_id = ? ", req.SpaceId)
	err = _db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = _db.Scopes(req.PageQuery()).Preload("User").Find(&list).Error
	return
}

func (r *memberRepository) Store(ctx context.Context, m *model.Member) error {
	return r.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "space_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"role"}),
	}).Create(&m).Error
}

func (r *memberRepository) GetByID(ctx context.Context, id int64) (*model.Member, error) {
	var m model.Member
	if err := r.DB(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *memberRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.DB(ctx).Delete(&model.Member{}, id).Error
}
