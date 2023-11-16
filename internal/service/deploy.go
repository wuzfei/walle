package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/wuzfei/go-helper/slices"
	"time"
	"yema.dev/api"
	"yema.dev/api/deploy"
	"yema.dev/internal/errcode"
	"yema.dev/internal/model"
	"yema.dev/internal/repository"
	"yema.dev/pkg/repo"
	"yema.dev/pkg/ssh"
)

type DeployService interface {
	List(ctx context.Context, req *deploy.ListReq) (total int64, res []*model.Task, err error)
	Create(ctx context.Context, req *deploy.CreateReq) error
	Delete(ctx context.Context, req *api.SpaceWithId) error
	Detail(ctx context.Context, req *api.SpaceWithId) (*model.Task, error)

	Audit(ctx context.Context, req *deploy.AuditReq) (err error)
	Release(ctx context.Context, req *api.SpaceWithId, userId int64) (err error)
	Stop(ctx context.Context, req *api.SpaceWithId) (err error)
	Rollback(ctx context.Context, req *api.SpaceWithId) (err error)
}

func NewDeployService(service *Service, deployRepo repository.DeployRepository, projectRepo repository.ProjectRepository, serverRepo repository.ServerRepository, ssh *ssh.Ssh, repo *repo.Repos) DeployService {
	return &deployService{
		deployRepo:  deployRepo,
		projectRepo: projectRepo,
		serverRepo:  serverRepo,
		Service:     service,
		ssh:         ssh,
		repo:        repo,
	}
}

type deployService struct {
	deployRepo  repository.DeployRepository
	serverRepo  repository.ServerRepository
	projectRepo repository.ProjectRepository
	ssh         *ssh.Ssh
	repo        *repo.Repos
	*Service
}

func (s *deployService) List(ctx context.Context, req *deploy.ListReq) (total int64, list []*model.Task, err error) {
	return s.deployRepo.List(ctx, req)
}

// Create 创建上线单
func (s *deployService) Create(ctx context.Context, req *deploy.CreateReq) error {
	project, err := s.projectRepo.GetByID(ctx, req.ProjectId)
	if err != nil {
		return err
	}
	if project.SpaceId != req.SpaceId {
		return errcode.ErrBadRequest
	}
	if !project.Status.IsEnable() || !project.Environment.Status.IsEnable() {
		return errors.New("该项目或者该环境暂停上线，请联系相关负责人")
	}
	serverIds := slices.Map(project.Servers, func(item model.Server, k int) int64 {
		return item.ID
	})
	m := &model.Task{
		Name:          req.Name,
		SpaceId:       req.SpaceId,
		UserId:        req.UserId,
		ProjectId:     project.ID,
		EnvironmentId: project.Environment.ID,
		Tag:           req.Tag,
		Branch:        req.Branch,
		CommitId:      req.CommitId,
	}
	m.Status = model.TaskStatusAudit
	if project.IsTaskAudit() {
		m.Status = model.TaskStatusWaiting
	}
	servers := make([]model.Server, 0)
	return s.tm.Transaction(ctx, func(ctx context.Context) (err error) {
		serverIds = slices.Intersect(serverIds, req.ServerIds)
		if len(serverIds) == 0 {
			return errcode.ErrBadRequest.Wrap(errors.New("服务器选择错误"))
		}
		servers, err = s.serverRepo.GetBySpaceAndIDs(ctx, req.SpaceId, serverIds)
		if err != nil {
			return err
		}
		if len(servers) == 0 {
			return errcode.ErrBadRequest.Wrap(errors.New("服务器选择错误"))
		}
		m.Servers = servers
		return s.deployRepo.Create(ctx, m)
	})
}

// Detail 上线单详情
func (s *deployService) Detail(ctx context.Context, req *api.SpaceWithId) (*model.Task, error) {
	m, err := s.deployRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if m.SpaceId != req.SpaceId {
		return nil, errcode.ErrBadRequest
	}
	return m, err
}

// Delete 删除
func (s *deployService) Delete(ctx context.Context, req *api.SpaceWithId) error {
	m, err := s.deployRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if m.SpaceId != req.SpaceId {
		return errcode.ErrBadRequest
	}
	return s.deployRepo.DeleteByID(ctx, req.ID)
}

// Audit 审核
func (s *deployService) Audit(ctx context.Context, req *deploy.AuditReq) (err error) {
	m, err := s.deployRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if m.SpaceId != req.SpaceId {
		return errcode.ErrBadRequest
	}
	if m.Status != model.TaskStatusWaiting {
		return errors.New("审核失败，该上线单并未处于待审核状态")
	}

	m.AuditUserId = req.AuditUserId
	if req.Audit {
		m.Status = model.TaskStatusAudit
	} else {
		m.Status = model.TaskStatusReject
	}
	m.AuditTime = sql.NullTime{Time: time.Now(), Valid: true}
	return s.deployRepo.Update(ctx, m, "status", "audit_user_id", "audit_time")
}

// Release 发布
func (s *deployService) Release(ctx context.Context, req *api.SpaceWithId, userId int64) (err error) {
	//上线单详情
	//taskDetail, err := srv.getTask(spaceAndId, "Project", "Environment", "Servers")
	//if err != nil {
	//	return
	//}
	//return srv.deploy.Start(taskDetail)
	return
}

// StopRelease 停止发布
func (s *deployService) Stop(ctx context.Context, req *api.SpaceWithId) (err error) {
	//上线单详情
	//taskDetail, err := srv.getTask(spaceAndId, "Project", "Environment", "Servers")
	//if err != nil {
	//	return
	//}
	//return srv.deploy.Stop(taskDetail.ID)
	return
}

// Rollback 回滚
func (s *deployService) Rollback(ctx context.Context, req *api.SpaceWithId) (err error) {
	//err = srv.db.First(&m, spaceId).Error
	return
}
