package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"yema.dev/internal/handler"
	"yema.dev/internal/model"
	"yema.dev/internal/repository"
)

func Permission(userRepo repository.UserRepository, role model.Role) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		log.Println("middleware Permission start")
		userId := handler.UserId(ctx)
		spaceId := handler.GetSpaceId(ctx)
		if !model.IsSuperUser(userId) {
			if spaceId == 0 {
				_ = ctx.AbortWithError(400, errors.New("未选择空间"))
				return
			}
			member, err := userRepo.GetMemberBySpaceAndUserId(ctx, spaceId, userId)
			if err != nil {
				_ = ctx.AbortWithError(400, errors.New("空间选择错误"))
				return
			}
			currRole := member.Role
			if model.Role(currRole).Level() < role.Level() {
				_ = ctx.AbortWithError(401, errors.New("你没有权限访问该空间，请联系相关负责人"))
				return
			}
		}
		ctx.Next()
		log.Println("middleware Permission end")
	}
}
