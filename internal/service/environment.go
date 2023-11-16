package service

import (
	"context"
	"errors"
	"yema.dev/api"
	"yema.dev/api/environment"
	"yema.dev/internal/errcode"
	"yema.dev/internal/model"
	"yema.dev/internal/repository"
)

type EnvironmentService interface {
	List(ctx context.Context, req *environment.ListReq) (total int64, res []*model.Environment, err error)
	Create(ctx context.Context, req *environment.CreateReq) error
	Update(ctx context.Context, req *environment.UpdateReq) error
	Delete(ctx context.Context, spaceWithId *api.SpaceWithId) error
}

func NewEnvironmentService(service *Service, environmentRepo repository.EnvironmentRepository) EnvironmentService {
	return &environmentService{
		environmentRepo: environmentRepo,
		Service:         service,
	}
}

type environmentService struct {
	environmentRepo repository.EnvironmentRepository
	*Service
}

func (s *environmentService) List(ctx context.Context, req *environment.ListReq) (total int64, list []*model.Environment, err error) {
	return s.environmentRepo.List(ctx, req)
}

func (s *environmentService) Create(ctx context.Context, req *environment.CreateReq) error {
	return s.environmentRepo.Create(ctx, &model.Environment{
		SpaceId:     req.SpaceId,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		Color:       req.Color,
	})
}

func (s *environmentService) Update(ctx context.Context, req *environment.UpdateReq) error {
	m, err := s.environmentRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if m.SpaceId != req.SpaceId {
		return errcode.ErrBadRequest
	}
	m.Name = req.Name
	m.Status = req.Status
	m.Description = req.Description
	m.Color = req.Color
	return s.environmentRepo.Update(ctx, m, req.Fields()...)
}

// Delete 环境下必须没有项目了才能删除
func (s *environmentService) Delete(ctx context.Context, spaceWithId *api.SpaceWithId) error {
	m, err := s.environmentRepo.GetByID(ctx, spaceWithId.ID)
	if err != nil {
		return err
	}
	if m.SpaceId != spaceWithId.SpaceId {
		return errcode.ErrBadRequest
	}
	if s.environmentRepo.GetProjectCount(ctx, spaceWithId.ID) > 0 {
		return errors.New("该环境还存在项目，不允许删除，如需要删除，请先删除该环境下所有项目")
	}
	return s.environmentRepo.DeleteByID(ctx, m.ID)
}
