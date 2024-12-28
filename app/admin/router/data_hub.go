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
	apiMkt := apis.DataHubMarketplace{}
	r := v1.Group("/data_hub").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionAction())
	{
		r.GET("/user", api.GetPageUser)
		r.GET("/user/point", api.GetPageUserPoint)
		r.GET("/marketplace/campaign", apiMkt.GetPageCampaign)
		r.GET("/marketplace/campaign/validation", apiMkt.GetCampaignValidation)
		r.GET("/marketplace/campaign/reward", apiMkt.GetPageReward)
		r.PUT("/marketplace/campaign/validation", apiMkt.UpdateCampaignValidation)
		r.POST("/marketplace/campaign", apiMkt.AddCampaign)
		r.PUT("/marketplace/campaign", apiMkt.UpdateCampaign)
		r.GET("/marketplace/campaign/:id", apiMkt.GetCampaignDetail)
		r.DELETE("/marketplace/campaign", apiMkt.DeleteCampaign)
		r.POST("/marketplace/campaign/upload", apiMkt.CampaignUpload)
		r.GET("/user/reward", api.GetAllReward)
		r.GET("/user/reward/export", api.GetAllRewardExport)
		r.GET("/marketplace/campaign/validation/dispute", apiMkt.GetCampaignDispute)
		r.PUT("/marketplace/campaign/validation/dispute", apiMkt.UpdateCampaignDispute)
		r.GET("/user/ambassador", api.GetPageAmbassadors)
		r.PUT("/user/ambassador", api.UpdateAmbassadors)
		r.GET("/user/ambassador/export", api.GetAmbassadorsExport)
		r.PUT("/marketplace/campaign/validation/malicious", apiMkt.UpdateCampaignValidationMalicious)
		r.GET("/marketplace/campaign/validation/summary", apiMkt.GetCampaignValidationSummary)
		r.GET("/marketplace/campaign/validation/download", apiMkt.GetCampaignValidationDownload)
	}
}
