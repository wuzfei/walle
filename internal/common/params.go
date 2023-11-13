package common

type StatisticsRes struct {
	Release  int `json:"release"`  //发布数总计
	Rollback int `json:"rollback"` //回滚数统计
	Projects int `json:"projects"` //项目数量
}
