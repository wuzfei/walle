package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"yema.dev/api"
	"yema.dev/api/user"
	"yema.dev/internal/errcode"
	"yema.dev/internal/model"
	"yema.dev/internal/repository"
	"yema.dev/pkg/jwt"
)

type UserService interface {
	Login(ctx context.Context, req *user.LoginReq) (*user.LoginRes, error)
	Logout(ctx context.Context, userId int64) (err error)
	RefreshToken(ctx context.Context, req *user.RefreshTokenReq) (res *user.LoginRes, err error)
	GetProfile(ctx context.Context, req *api.SpaceWithId) (*user.ProfileRes, error)

	List(ctx context.Context, req *user.ListReq) (total int64, res []*model.User, err error)
	Create(ctx context.Context, req *user.CreateReq) error
	Update(ctx context.Context, req *user.UpdateReq) error
	Delete(ctx context.Context, id int64) error
}

func NewUserService(service *Service, userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
		Service:  service,
	}
}

type userService struct {
	userRepo repository.UserRepository
	*Service
}

func (s *userService) Login(ctx context.Context, req *user.LoginReq) (*user.LoginRes, error) {\
	m, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil  {
		return nil, err
	}
	if m == nil ||  m.Status.IsDisable() {
		return nil, errcode.ErrUnauthorized
	}
	if bcrypt.CompareHashAndPassword(m.Password, []byte(req.Password)) != nil {
		return nil, errcode.ErrUnauthorized
	}
	//生成token
	res := user.LoginRes{}
	res.Token, res.TokenExpire, err = s.jwt.CreateToken(jwt.TokenPayload{
		UserId:   m.ID,
		Email:    m.Email,
		Username: m.Username,
	})
	if err != nil {
		return nil, err
	}
	res.UserId = m.ID
	//记住登陆
	if req.Remember {
		res.RefreshToken, res.RefreshTokenExpire, err = s.jwt.CreateRefreshToken(jwt.TokenPayload{
			UserId:    m.ID,
			Email:     m.Email,
			Username:  m.Username,
			IsRefresh: true,
		})
		if err != nil {
			return nil, err
		}
		m.RememberToken = res.RefreshToken
		if err = s.userRepo.UpdateFields(ctx, m, "remember_token"); err != nil {
			return nil, err
		}
	}
	return &res, nil
}

func (s *userService) GetProfile(ctx context.Context, req *api.SpaceWithId) ( *user.ProfileRes,  error) {
	m, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	role := ""
	if m.IsSuperUser() {
		role = string(model.RoleSuper)
	}
	currentSpaceId := req.SpaceId
	//获取所属空间
	spaceItems, err := s.spacesItems(ctx, m.ID)
	if err != nil {
		return nil, err
	}
	//根据当前空间，获取当前空间id和角色
	currSpaceItem := spaceItems.Default(currentSpaceId)
	if currSpaceItem != nil {
		currentSpaceId = currSpaceItem.SpaceId
		if role == "" {
			role = currSpaceItem.Role
		}
	} else {
		currentSpaceId = 0
	}

	return &user.ProfileRes{
		UserID:         m.ID,
		Email:          m.Email,
		Username:       m.Username,
		Role:           role,
		Status:         m.Status,
		CurrentSpaceId: currentSpaceId,
		Spaces:         spaceItems,
	}, nil
}

// RefreshToken 刷新token
func (s *userService) RefreshToken(ctx context.Context, req *user.RefreshTokenReq) (res *user.LoginRes, err error) {
	jwtClaims, err := s.jwt.ValidateToken(req.RefreshToken)
	if err != nil {
		return
	}
	m, err := s.userRepo.GetByID(ctx, jwtClaims.UserId)
	if err != nil {
		return
	}
	if m.Status.IsDisable() {
		return nil, errcode.ErrUnauthorized
	}
	if m.RememberToken != req.RefreshToken {
		return nil, errors.New("refresh token 错误")
	}

	res = &user.LoginRes{}
	res.Token, res.TokenExpire, err = s.jwt.CreateToken(jwt.TokenPayload{
		UserId:   m.ID,
		Email:    m.Email,
		Username: m.Username,
	})
	if err != nil {
		return
	}
	res.UserId = m.ID
	res.RefreshToken, res.RefreshTokenExpire, err = s.jwt.CreateRefreshToken(jwt.TokenPayload{
		UserId:    m.ID,
		Email:     m.Email,
		Username:  m.Username,
		IsRefresh: true,
	})
	if err != nil {
		return nil, err
	}
	m.RememberToken = res.RefreshToken
	if err = s.userRepo.UpdateFields(ctx, m, "remember_token"); err != nil {
		return nil, err
	}
	return
}

// Logout 退出
func (s *userService) Logout(ctx context.Context, userId int64) (err error) {
	m, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return
	}
	m.RememberToken = ""
	if err = s.userRepo.UpdateFields(ctx, m, "remember_token"); err != nil {
		return err
	}
	return nil
}

func (s *userService) List(ctx context.Context, req *user.ListReq) (total int64, res []*model.User, err error) {
	return s.userRepo.List(ctx, req.Keyword, req.PageQuery())
}

func (s *userService) Create(ctx context.Context, req *user.CreateReq) error {
	m, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if m.ID != 0 {
		return errors.New("该用户email已存在")
	}
	m.Username = req.Username
	m.Email = req.Email
	m.Status = req.Status

	_pwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.Password = _pwd
	return s.userRepo.Create(ctx, m)
}

func (s *userService) Update(ctx context.Context, req *user.UpdateReq) error {
	m, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if m.ID == 0 {
		return errors.New("用户不存在")
	}
	m.ID = req.ID
	m.Username = req.Username
	m.Email = req.Email
	m.Status = req.Status
	if req.Password != "" {
		m.Password, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
	}
	return s.userRepo.Update(ctx, m)
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	if model.IsSuperUser(id) {
		return errors.New("超级管理员不允许删除")
	}
	return s.userRepo.DeleteByID(ctx, id)
}

// spacesItems 获取用户所属的所有空间
func (s *userService) spacesItems(ctx context.Context, userId int64) (user.SpaceItems, error) {
	spaceItems := make(user.SpaceItems, 0)
	if !model.IsSuperUser(userId) {
		res, err := s.userRepo.GetMemberSpaces(ctx, userId)
		if err != nil {
			return spaceItems, err
		}
		for _, item := range res {
			if item.SpaceId != 0 {
				spaceItems = append(spaceItems, &user.SpaceItem{
					SpaceId:   item.Space.ID,
					SpaceName: item.Space.Name,
					Status:    item.Space.Status,
					Role:      item.Role,
				})
			}
		}
	} else {
		//超级管理员的处理
		//var res []*model.Space
		//err := srv.db.Find(&res).Error
		//if err != nil {
		//	return spaceItems, err
		//}
		//spaceItems = slices.Map(res, func(item *model.Space, k int) *SpaceItem {
		//	return &SpaceItem{
		//		SpaceId:   item.ID,
		//		SpaceName: item.Name,
		//		Status:    item.Status,
		//		Role:      string(model.RoleSuper),
		//	}
		//})
	}
	return spaceItems, nil
}
