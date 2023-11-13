package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"yema.dev/internal/errcode"
	"yema.dev/internal/model"
)

type UserRepository interface {
	List(ctx context.Context, kw string, scopesFn ...func(*gorm.DB) *gorm.DB) (total int64, res []*model.User, err error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	UpdateFields(ctx context.Context, user *model.User, fields ...string) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	DeleteByID(ctx context.Context, id int64) error
	GetMemberSpaces(ctx context.Context, id int64) ([]*model.Member, error)
}

func NewUserRepository(r *Repository) UserRepository {
	return &userRepository{
		Repository: r,
	}
}

type userRepository struct {
	*Repository
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.DB(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.DB(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) UpdateFields(ctx context.Context, user *model.User, fields ...string) error {
	if err := r.DB(ctx).Select(fields).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, userId int64) (*model.User, error) {
	var user model.User
	if err := r.DB(ctx).Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.DB(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) DeleteByID(ctx context.Context, userId int64) error {
	return r.DB(ctx).Delete(&model.User{}, userId).Error
}

// GetMemberSpaces 获取一个用户所有空间信息
func (r *userRepository) GetMemberSpaces(ctx context.Context, userId int64) (res []*model.Member, err error) {
	err = r.DB(ctx).Where("user_id = ?", userId).Preload("Space").Find(&res).Error
	return
}

// List 获取列表
func (r *userRepository) List(ctx context.Context, kw string, scopesFn ...func(*gorm.DB) *gorm.DB) (total int64, res []*model.User, err error) {
	db := r.DB(ctx).Model(&model.User{})
	if kw != "" {
		_k := "%" + kw + "%"
		db.Where("username like ? or email like ?", _k, _k)
	}
	err = db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = db.Scopes(scopesFn...).Find(&res).Error
	return
}
