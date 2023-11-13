package service

import (
	"context"
	"errors"
	"yema.dev/api/space"
	"yema.dev/internal/model"
	"yema.dev/internal/repository"
)

type SpaceService interface {
	List(ctx context.Context, req *space.ListReq) (total int64, res []*model.Space, err error)
	Create(ctx context.Context, req *space.CreateReq) error
	Update(ctx context.Context, req *space.UpdateReq) error
	Delete(ctx context.Context, id int64) error
}

func NewSpaceService(service *Service, spaceRepo repository.SpaceRepository) SpaceService {
	return &spaceService{
		spaceRepo: spaceRepo,
		Service:   service,
	}
}

type spaceService struct {
	spaceRepo repository.SpaceRepository
	*Service
}

func (s *spaceService) List(ctx context.Context, req *space.ListReq) (total int64, res []*model.Space, err error) {
	return s.spaceRepo.List(ctx, req)
}

func (s *spaceService) Create(ctx context.Context, req *space.CreateReq) error {
	return s.spaceRepo.Create(ctx, &model.Space{
		UserId: req.UserId,
		Name:   req.Name,
		Status: req.Status,
	})
}

func (s *spaceService) Update(ctx context.Context, req *space.UpdateReq) error {
	return s.spaceRepo.UpdateFields(ctx, &model.Space{
		ID:     req.ID,
		Name:   req.Name,
		Status: req.Status,
		UserId: req.UserId,
	}, req.Fields()...)
}

func (s *spaceService) Delete(ctx context.Context, id int64) error {
	if s.spaceRepo.GetProjectCount(ctx, id) > 0 {
		return errors.New("该空间存在项目，不允许删除，如需要删除，先删除该空间下所有项目")
	}
	return s.spaceRepo.DeleteByID(ctx, id)
}
