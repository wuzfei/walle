package model

import (
	"gorm.io/gorm"
	"time"
	"yema.dev/app/model/field"
)

const (
	ProjectIsInclude = 0
	ProjectIsExclude = 1

	ProjectTaskAuditEnable  = 1
	ProjectTaskAuditDisable = 0
)

type Project struct {
	ID            int64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
	SpaceId       int64 `gorm:"column:space_id;index;notNull;comment:所属空间" json:"space_id"`
	EnvironmentId int64 `gorm:"column:environment_id;notNull;comment:所属环境" json:"environment_id"`
	UserId        int64 `gorm:"column:user_id;notNull;comment:所属用户" json:"user_id"`

	Name         string `gorm:"column:name;size:100;notNull;comment:名称" json:"name"`
	RepoUrl      string `gorm:"column:repo_url;size:500;notNull;comment:仓库地址" json:"repo_url"`
	RepoMode     string `gorm:"column:repo_mode;size:10;notNull;default:tag;comment:分支类型" json:"repo_mode"` // tag/branch
	RepoUsername string `gorm:"column:repo_username;size:100;notNull;default:'';comment:仓库用户" json:"repo_username"`
	RepoPassword string `gorm:"column:repo_password;size:100;notNull;default:'';comment:仓库密码" json:"repo_password"`
	RepoType     string `gorm:"column:repo_type;size:20;notNull;default:git;comment:仓库类型" json:"repo_type"`

	Excludes  string `gorm:"column:excludes;size:1000;notNull;default:'';comment:包含或者去除的文件列表" json:"excludes"` //包含或者去除的文件
	IsInclude int8   `gorm:"column:is_include;notNull;default:0;comment:1去除0包含" json:"is_include"`             //是包含还是去除

	TaskVars    string `gorm:"column:task_vars;size:1000;notNull;default:'';comment:全局环境变量" json:"task_vars"` //全局变量
	PrevDeploy  string `gorm:"column:prev_deploy;size:1000;notNull;default:'';comment:编译前操作命令" json:"prev_deploy"`
	PostDeploy  string `gorm:"column:post_deploy;size:1000;notNull;default:'';comment:编译后操作命令" json:"post_deploy"`
	PrevRelease string `gorm:"column:prev_release;size:1000;notNull;default:'';comment:发布前操作命令" json:"prev_release"`
	PostRelease string `gorm:"column:post_release;size:1000;notNull;default:'';comment:发布后操作命令" json:"post_release"`

	TargetRoot     string `gorm:"column:target_root;size:500;notNull;default:'';comment:目标路径" json:"target_root"` //目标路径
	TargetReleases string `gorm:"column:target_releases;size:500;notNull;default:'';comment:目标代码路径" json:"target_releases"`
	KeepVersionNum int    `gorm:"column:keep_version_num;notNull;default:5;comment:保留版本数量" json:"keep_version_num"` //保留版本数量
	TaskAudit      int8   `gorm:"column:task_audit;notNull;default:1;comment:上线单是否开启审核" json:"task_audit"`          //上线单是否开启审核
	Description    string `gorm:"column:description;size:500;notNull;default:'';comment:简介说明" json:"description"`

	Master     string `gorm:"column:master" json:"master"`
	Version    string `gorm:"column:version;size:100;notNull;default:'';comment:版本号" json:"version"`
	NoticeType string `gorm:"column:notice_type" json:"notice_type"`
	NoticeHook string `gorm:"column:notice_hook" json:"notice_hook"`

	Status field.Status `gorm:"column:status;size:1;notNull;default:0;comment:状态" json:"status"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;notNull" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;notNull" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	Space       *Space       `json:"space,omitempty"`
	Environment *Environment `json:"environment,omitempty"`
	Servers     []Server     `gorm:"many2many:project_server" json:"servers,omitempty"`
}

func (p *Project) IsTaskAudit() bool {
	return p.TaskAudit == ProjectTaskAuditEnable
}
