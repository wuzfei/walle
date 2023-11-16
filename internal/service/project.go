package service

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"strconv"
	"time"
	"yema.dev/api"
	"yema.dev/api/project"
	"yema.dev/internal/errcode"
	"yema.dev/internal/model"
	"yema.dev/internal/model/field"
	"yema.dev/internal/repository"
	"yema.dev/pkg/repo"
	"yema.dev/pkg/ssh"
)

type ProjectService interface {
	List(ctx context.Context, req *project.ListReq) (total int64, res []*model.Project, err error)
	Create(ctx context.Context, req *project.CreateReq) error
	Update(ctx context.Context, req *project.UpdateReq) error
	Delete(ctx context.Context, req *api.SpaceWithId) error
	Detail(ctx context.Context, req *api.SpaceWithId) (*model.Project, error)

	DetectionWs(ctx context.Context, wsConn *websocket.Conn, spaceWithId *api.SpaceWithId) (err error)

	GetBranches(ctx context.Context, req *api.SpaceWithId) (res []repo.Branch, err error)
	GetTags(ctx context.Context, req *api.SpaceWithId) (res []repo.Tag, err error)
	GetCommits(ctx context.Context, req *api.SpaceWithId, branch string) (res []repo.Commit, err error)
}

func NewProjectService(service *Service, projectRepo repository.ProjectRepository, serverRepo repository.ServerRepository, ssh *ssh.Ssh, repo *repo.Repos) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		serverRepo:  serverRepo,
		Service:     service,
		ssh:         ssh,
		repo:        repo,
	}
}

type projectService struct {
	projectRepo repository.ProjectRepository
	serverRepo  repository.ServerRepository
	ssh         *ssh.Ssh
	repo        *repo.Repos
	*Service

	detectionTimeout time.Duration //检测项目时的超时时间
}

func (s *projectService) List(ctx context.Context, req *project.ListReq) (total int64, list []*model.Project, err error) {
	return s.projectRepo.List(ctx, req)
}

func (s *projectService) Create(ctx context.Context, req *project.CreateReq) error {
	m := &model.Project{
		SpaceId: req.SpaceId,

		Name:          req.Name,
		EnvironmentId: req.EnvironmentId,
		RepoUrl:       req.RepoUrl,
		RepoMode:      req.RepoMode,
		RepoType:      req.RepoType,
		TaskAudit:     req.TaskAudit,
		Description:   req.Description,

		TargetRoot:     req.TargetRoot,
		TargetReleases: req.TargetReleases,
		KeepVersionNum: req.KeepVersionNum,

		Excludes:    req.Excludes,
		IsInclude:   req.IsInclude,
		TaskVars:    req.TaskVars,
		PrevDeploy:  req.PrevRelease,
		PostDeploy:  req.PostDeploy,
		PrevRelease: req.PrevRelease,
		PostRelease: req.PostRelease,
		Status:      field.StatusEnable,
	}
	servers := make([]model.Server, 0)
	return s.tm.Transaction(ctx, func(ctx context.Context) (err error) {
		servers, err = s.serverRepo.GetBySpaceAndIDs(ctx, req.SpaceId, req.ServerIds)
		if err != nil {
			return err
		}
		m.Servers = servers
		return s.projectRepo.Create(ctx, m)
	})
}

func (s *projectService) Update(ctx context.Context, req *project.UpdateReq) error {
	m, err := s.projectRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if m.SpaceId != req.SpaceId {
		return errcode.ErrBadRequest
	}

	m = &model.Project{
		ID:      req.ID,
		SpaceId: req.SpaceId,

		Name:          req.Name,
		EnvironmentId: req.EnvironmentId,
		RepoUrl:       req.RepoUrl,
		RepoMode:      req.RepoMode,
		RepoType:      req.RepoType,
		TaskAudit:     req.TaskAudit,
		Description:   req.Description,

		TargetRoot:     req.TargetRoot,
		TargetReleases: req.TargetReleases,
		KeepVersionNum: req.KeepVersionNum,

		Excludes:    req.Excludes,
		IsInclude:   req.IsInclude,
		TaskVars:    req.TaskVars,
		PrevDeploy:  req.PrevRelease,
		PostDeploy:  req.PostDeploy,
		PrevRelease: req.PrevRelease,
		PostRelease: req.PostRelease,
	}
	return s.tm.Transaction(ctx, func(ctx context.Context) error {
		servers, err := s.serverRepo.GetBySpaceAndIDs(ctx, req.SpaceId, req.ServerIds)
		if err != nil {
			return err
		}

		//晴空关联
		err = s.projectRepo.ClearServers(ctx, req.ID)
		if err != nil {
			return err
		}
		//更新关联
		m.Servers = servers
		_fields := req.Fields()
		_fields = append(_fields, "Servers")
		return s.projectRepo.Update(ctx, m, _fields...)
	})
}

func (s *projectService) Detail(ctx context.Context, req *api.SpaceWithId) (*model.Project, error) {
	m, err := s.projectRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if m.SpaceId != req.SpaceId {
		return nil, errcode.ErrBadRequest
	}
	return m, err
}

// Delete 删除
func (s *projectService) Delete(ctx context.Context, req *api.SpaceWithId) error {
	return s.tm.Transaction(ctx, func(ctx context.Context) error {
		m, err := s.projectRepo.GetByID(ctx, req.ID)
		if err != nil {
			return err
		}
		if m.SpaceId != req.SpaceId {
			return errcode.ErrBadRequest
		}
		if err = s.projectRepo.ClearServers(ctx, req.ID); err != nil {
			return err
		}
		return s.projectRepo.DeleteByID(ctx, m.ID)
	})
}

func (s *projectService) GetBranches(ctx context.Context, req *api.SpaceWithId) (res []repo.Branch, err error) {
	var rep repo.Repo
	rep, err = s.getRepoBySpaceWithId(ctx, req)
	if err != nil {
		return
	}
	return rep.Branches()
}

func (s *projectService) GetTags(ctx context.Context, req *api.SpaceWithId) (res []repo.Tag, err error) {
	var rep repo.Repo
	rep, err = s.getRepoBySpaceWithId(ctx, req)
	if err != nil {
		return
	}
	return rep.Tags()
}

func (s *projectService) GetCommits(ctx context.Context, req *api.SpaceWithId, branch string) (res []repo.Commit, err error) {
	var rep repo.Repo
	rep, err = s.getRepoBySpaceWithId(ctx, req)
	if err != nil {
		return
	}
	return rep.Commits(branch)
}

func (s *projectService) getRepoBySpaceWithId(ctx context.Context, req *api.SpaceWithId) (rep repo.Repo, err error) {
	projectModel, err := s.Detail(ctx, req)
	if err != nil {
		return nil, err
	}
	if !projectModel.Status.IsEnable() {
		return nil, errors.New("该项目已经禁用")
	}
	return s.repo.New(repo.TypeRepo(projectModel.RepoType), projectModel.RepoUrl, strconv.Itoa(int(projectModel.ID)))
}
