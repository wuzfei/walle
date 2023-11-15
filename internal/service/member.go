package service

import (
	"context"
	"github.com/wuzfei/go-helper/slices"
	"github.com/zeebo/errs"
	"yema.dev/api"
	"yema.dev/api/member"
	"yema.dev/internal/errcode"
	"yema.dev/internal/model"
	"yema.dev/internal/repository"
	"yema.dev/pkg/ssh"
)

type MemberService interface {
	List(ctx context.Context, req *member.ListReq) (total int64, res []*member.ListItem, err error)
	Store(ctx context.Context, req *member.StoreReq) error
	Delete(ctx context.Context, req *api.SpaceWithId) error
}

func NewMemberService(service *Service, memberRepo repository.MemberRepository, userRepo repository.UserRepository) MemberService {
	return &memberService{
		memberRepo: memberRepo,
		userRepo:   userRepo,
		Service:    service,
	}
}

type memberService struct {
	memberRepo repository.MemberRepository
	userRepo   repository.UserRepository
	ssh        *ssh.Ssh
	*Service
}

func (s *memberService) List(ctx context.Context, req *member.ListReq) (total int64, res []*member.ListItem, err error) {
	var members []*model.Member
	total, members, err = s.memberRepo.List(ctx, req)
	if err != nil {
		return
	}
	res = slices.Map(members, func(item *model.Member, k int) *member.ListItem {
		return &member.ListItem{
			SpaceId:   item.SpaceId,
			UserId:    item.UserId,
			Username:  item.User.Username,
			Email:     item.User.Email,
			Role:      item.Role,
			Status:    item.User.Status,
			CreatedAt: item.CreatedAt,
		}
	})
	return
}

func (s *memberService) Store(ctx context.Context, req *member.StoreReq) error {
	user, err := s.userRepo.GetByID(ctx, req.UserId)
	if err != nil {
		return err
	}
	if !user.Status.IsEnable() {
		return errs.New("该用户已被禁用")
	}
	return s.memberRepo.Store(ctx, &model.Member{
		SpaceId: req.SpaceId,
		UserId:  req.UserId,
		Role:    req.Role,
	})
}

// Delete 删除
func (s *memberService) Delete(ctx context.Context, req *api.SpaceWithId) error {
	m, err := s.memberRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if m.SpaceId != req.SpaceId {
		return errcode.ErrBadRequest
	}
	return s.memberRepo.DeleteByID(ctx, req.ID)
}
