package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerDataHubUserRouter)
}

// 需认证的路由代码
func registerDataHubUserRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.DataHubUser{}
	r := v1.Group("/data_hub").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionAction())
	{
		r.GET("/user", api.GetPageUser)
		r.GET("/user/point", api.GetPageUserPoint)
	}
}
