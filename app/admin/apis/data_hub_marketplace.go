package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"strconv"
)

type DataHubMarketplace struct {
	api.Api
}

// GetPageCampaign
// @Summary 列表DataHub所有Campaign信息数据
// @Description 获取JSON
// @Tags DataHub
// @Success 200 {string} {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/data_hub/marketplace/campaign [get]
// @Security Bearer
func (e DataHubMarketplace) GetPageCampaign(c *gin.Context) {
	s := service.DataHubMarketplace{}
	req := dto.DataHubMarketplaceGetPageCampaignReq{}

	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	list := make([]models.AITask, 0)
	var count int64

	err = s.GetPageCampaign(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	for i, item := range list {
		list[i].USDTReward = item.USDTReward / 100
		list[i].PointReward = item.PointReward / 100
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetCampaignValidation Get
// @Summary 获取Campaign信息数据
// @Description 获取JSON
// @Tags DataHub
// @Param task_id query string false "task_id"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/data_hub/marketplace/campaign/validation [get]
// @Security Bearer
func (e DataHubMarketplace) GetCampaignValidation(c *gin.Context) {
	s := service.DataHubMarketplace{}
	req := dto.DataHubMarketplaceGetCampaignValidationReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	list := make([]models.AITaskUploadRecord, 0)
	var count int64

	err = s.GetCampaignValidation(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	for i, item := range list {
		var files []models.AITaskUploadedFile
		e.Orm.Find(&files, "upload_record = ?", item.ID)
		list[i].Items = files
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetPageReward Get
// @Summary 获取Campaign reward信息数据
// @Description 获取JSON
// @Tags DataHub
// @Param task_id query string false "task_id"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/data_hub/marketplace/campaign/reward [get]
// @Security Bearer
func (e DataHubMarketplace) GetPageReward(c *gin.Context) {
	s := service.DataHubMarketplace{}
	req := dto.DataHubMarketplaceGetPageRewardReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	list := make([]models.MarketPlaceRewardItem, 0)
	var count int64

	err = s.GetCampaignReward(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	for i, item := range list {
		usdtIntVal, err := strconv.Atoi(item.USDTReward)
		if err != nil {
			usdtIntVal = 0
		}
		usdtFloatVal := float64(usdtIntVal) / 100
		usdt := fmt.Sprintf("%.2f", usdtFloatVal)
		pointIntVal, err := strconv.Atoi(item.PointReward)
		if err != nil {
			pointIntVal = 0
		}
		pointFloatVal := float64(pointIntVal) / 100
		point := fmt.Sprintf("%.0f", pointFloatVal)
		list[i].USDTReward = usdt
		list[i].PointReward = point
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// UpdateCampaignValidation Update 修改validation
// @Summary 修改validation
// @Description 获取JSON
// @Tags DataHub
// @Accept  application/json
// @Product application/json
// @Param data body dto.MarketplaceValidationUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/data_hub/marketplace/campaign/validation [put]
// @Security Bearer
func (e DataHubMarketplace) UpdateCampaignValidation(c *gin.Context) {
	s := service.DataHubMarketplace{}
	req := dto.MarketplaceValidationUpdateReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}
	p := actions.GetPermissionFromContext(c)
	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(req.GetId(), "更新成功")
}
