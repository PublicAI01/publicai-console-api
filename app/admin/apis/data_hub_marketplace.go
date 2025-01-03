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
	"time"
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

	//list := make([]models.AITask, 0)
	//var count int64
	//
	//err = s.GetPageCampaign(&req, p, &list, &count)
	//if err != nil {
	//	e.Error(500, err, "查询失败")
	//	return
	//}
	//for i, item := range list {
	//	list[i].USDTReward = item.USDTReward / 100
	//	list[i].PointReward = item.PointReward / 100
	//}
	list := make([]models.AITaskVariants, 0)
	var count int64

	err = s.GetPageCampaignVariants(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetCampaignValidation Get
// @Summary 获取Campaign信息数据
// @Description 获取JSON
// @Tags DataHub
// @Param task_id query string false "task_id"
// @Param start_time query string false "start_time"
// @Param end_time query string false "end_time"
// @Param uid query string false "uid"
// @Param status query string false "status"
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

	list := make([]models.AITaskShowRecordItem, 0)
	var count int64
	if req.StartTime != "" && req.EndTime != "" {
		startTimeStamp, err := strconv.ParseInt(req.StartTime, 10, 64)
		if err != nil {
			e.Logger.Error(err)
			e.Error(500, err, err.Error())
			return
		}
		endTimeStamp, err := strconv.ParseInt(req.EndTime, 10, 64)
		if err != nil {
			e.Logger.Error(err)
			e.Error(500, err, err.Error())
			return
		}
		startTime := time.Unix(startTimeStamp, 0)
		endTime := time.Unix(endTimeStamp, 0)
		req.StartTime = startTime.Format("2006-01-02 15:04:05")
		req.EndTime = endTime.Format("2006-01-02 15:04:05")
	}
	err = s.GetCampaignValidation(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}

	for i, item := range list {
		var files []models.AITaskUploadedFile
		e.Orm.Find(&files, "upload_record = ?", item.ID)
		fileItems := make([]models.FileItem, len(files))
		for j, file := range files {
			fileItems[j] = models.FileItem{
				ID:    int(file.ID),
				Link:  file.Link,
				VPass: file.VPass,
				APass: file.APass,
			}
		}
		list[i].Items = fileItems
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

// GetCampaignDispute Get
// @Summary 获取Campaign争议的题目信息数据
// @Description 获取JSON
// @Tags DataHub
// @Param task_id query string false "task_id"
// @Param uid query string false "uid"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/data_hub/marketplace/campaign/validation/dispute [get]
// @Security Bearer
func (e DataHubMarketplace) GetCampaignDispute(c *gin.Context) {
	s := service.DataHubMarketplace{}
	req := dto.DataHubMarketplaceGetCampaignDisputeReq{}
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

	list := make([]models.AITaskShowDisputeRecordItem, 0)
	var count int64

	err = s.GetCampaignDispute(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}

	for i, item := range list {
		var files []models.AITaskUploadedFile
		e.Orm.Find(&files, "upload_record = ?", item.ID)
		fileItems := make([]models.FileDisputeItem, len(files))
		for j, file := range files {
			fileItems[j] = models.FileDisputeItem{
				ID:   int(file.ID),
				Link: file.Link,
				VAye: file.VAye,
				VNay: file.VNay,
			}
		}
		list[i].Items = fileItems
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// UpdateCampaignDispute Update 修改dispute
// @Summary 修改dispute
// @Description 获取JSON
// @Tags DataHub
// @Accept  application/json
// @Product application/json
// @Param data body dto.MarketplaceDisputeUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/data_hub/marketplace/campaign/validation/dispute [put]
// @Security Bearer
func (e DataHubMarketplace) UpdateCampaignDispute(c *gin.Context) {
	s := service.DataHubMarketplace{}
	req := dto.MarketplaceDisputeUpdateReq{}
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
	err = s.UpdateDispute(&req, p)
	if err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(req.GetId(), "更新成功")
}

// UpdateCampaignValidationMalicious Update 修改或撤销validation的Malicious状态
// @Summary 修改或撤销validation的Malicious状态
// @Description 获取JSON
// @Tags DataHub
// @Accept  application/json
// @Product application/json
// @Param data body dto.MarketplaceValidationMaliciousUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/data_hub/marketplace/campaign/validation/malicious [put]
// @Security Bearer
func (e DataHubMarketplace) UpdateCampaignValidationMalicious(c *gin.Context) {
	s := service.DataHubMarketplace{}
	req := dto.MarketplaceValidationMaliciousUpdateReq{}
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
	err = s.UpdateMalicious(&req, p)
	if err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(req.GetId(), "更新成功")
}

// GetCampaignValidationSummary Get
// @Summary 获取Campaign Validation Summary信息数据
// @Description 获取JSON
// @Tags DataHub
// @Param task_id query string false "task_id"
// @Param start_time query string false "start_time"
// @Param end_time query string false "end_time"
// @Param uid query string false "uid"
// @Param status query string false "status"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/data_hub/marketplace/campaign/validation/summary [get]
// @Security Bearer
func (e DataHubMarketplace) GetCampaignValidationSummary(c *gin.Context) {
	s := service.DataHubMarketplace{}
	req := dto.GetCampaignValidationSummaryReq{}
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
	startTimeStamp, err := strconv.ParseInt(req.StartTime, 10, 64)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	endTimeStamp, err := strconv.ParseInt(req.EndTime, 10, 64)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	startTime := time.Unix(startTimeStamp, 0)
	endTime := time.Unix(endTimeStamp, 0)
	req.StartTime = startTime.Format("2006-01-02 15:04:05")
	req.EndTime = endTime.Format("2006-01-02 15:04:05")
	var object models.ValidationSummary
	err = s.GetValidationSummary(&req, p, &object).Error
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(object, "查询成功")
}

// GetCampaignValidationDownload Get
// @Summary 获取待下载的图片信息
// @Description 获取JSON
// @Tags DataHub
// @Param task_id query string false "task_id"
// @Param start_time query string false "start_time"
// @Param end_time query string false "end_time"
// @Param uid query string false "uid"
// @Param status query string false "status"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/data_hub/marketplace/campaign/validation/download [get]
// @Security Bearer
func (e DataHubMarketplace) GetCampaignValidationDownload(c *gin.Context) {
	s := service.DataHubMarketplace{}
	req := dto.GetCampaignValidationSummaryReq{}
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
	startTimeStamp, err := strconv.ParseInt(req.StartTime, 10, 64)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	endTimeStamp, err := strconv.ParseInt(req.EndTime, 10, 64)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	startTime := time.Unix(startTimeStamp, 0)
	endTime := time.Unix(endTimeStamp, 0)
	req.StartTime = startTime.Format("2006-01-02 15:04:05")
	req.EndTime = endTime.Format("2006-01-02 15:04:05")
	list := make([]models.AITaskShowRecordItem, 0)
	err = s.DownloadValidation(&req, p, &list)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	for i, item := range list {
		var files []models.AITaskUploadedFile
		e.Orm.Find(&files, "upload_record = ?", item.ID)
		fileItems := make([]models.FileItem, len(files))
		for j, file := range files {
			fileItems[j] = models.FileItem{
				ID:    int(file.ID),
				Link:  file.Link,
				VPass: file.VPass,
				APass: file.APass,
			}
		}
		list[i].Items = fileItems
	}
	e.OK(list, "查询成功")
}

// AddCampaign 新增campaign
// @Summary 新增campaign
// @Description 获取JSON
// @Tags DataHub
// @Accept  application/json
// @Product application/json
// @Param data body dto.AddCampaignReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/data_hub/marketplace/campaign [post]
// @Security Bearer
func (e DataHubMarketplace) AddCampaign(c *gin.Context) {
	s := service.DataHubMarketplace{}
	req := dto.AddCampaignReq{}
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
	err = s.AddCampaign(&req, p)
	if err != nil {
		e.Error(500, err, "Add campaign failed")
		return
	}
	e.OK(0, "Add campaign success")
}

// UpdateCampaign 更新campaign
// @Summary 更新campaign
// @Description 获取JSON
// @Tags DataHub
// @Accept  application/json
// @Product application/json
// @Param data body dto.UpdateCampaignReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/data_hub/marketplace/campaign [put]
// @Security Bearer
func (e DataHubMarketplace) UpdateCampaign(c *gin.Context) {
	s := service.DataHubMarketplace{}
	req := dto.UpdateCampaignReq{}
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
	err = s.UpdateCampaign(&req, p)
	if err != nil {
		e.Error(500, err, "Update campaign failed")
		return
	}
	e.OK(0, "Update campaign success")
}

// DeleteCampaign 删除campaign
// @Summary 删除campaign
// @Description 获取JSON
// @Tags DataHub
// @Accept  application/json
// @Product application/json
// @Param data body dto.DeleteCampaignReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/data_hub/marketplace/campaign [delete]
// @Security Bearer
func (e DataHubMarketplace) DeleteCampaign(c *gin.Context) {
	s := service.DataHubMarketplace{}
	req := dto.DeleteCampaignReq{}
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
	err = s.DeleteCampaign(&req, p)
	if err != nil {
		e.Error(500, err, "Delete campaign failed")
		return
	}
	e.OK(0, "Delete campaign success")
}

// CampaignUpload 上传图片
// @Summary 上传图片
// @Description 获取JSON
// @Tags DataHub
// @Param data body dto.CampaignUploadReq true "body"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/data_hub/marketplace/campaign/upload [post]
// @Security Bearer
func (e DataHubMarketplace) CampaignUpload(c *gin.Context) {
	s := service.DataHubMarketplace{}
	req := dto.CampaignUploadReq{}
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
	if err = c.ShouldBind(&req); err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//数据权限检查
	p := actions.GetPermissionFromContext(c)
	var object dto.CampaignUploadResponse
	err = s.CampaignUpload(&req, p, &object)
	if err != nil {
		e.Error(500, err, "Upload failed")
		return
	}
	e.OK(object, "Upload success")
}

// GetCampaignDetail 获取某个campaign信息
// @Summary 获取某个campaign信息
// @Description 获取某个campaign信息
// @Tags DataHub
// @Param id path string false "id"
// @Success 200 {object} response.Response{data=dto.CampaignDetailResponse} "{"code": 200, "data": [...]}"
// @Router /api/v1/data_hub/marketplace/campaign/{id} [get]
// @Security Bearer
func (e DataHubMarketplace) GetCampaignDetail(c *gin.Context) {
	req := dto.CampaignDetailReq{}
	s := service.DataHubMarketplace{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	p := actions.GetPermissionFromContext(c)
	var object dto.CampaignDetailResponse
	err = s.CampaignDetail(&req, p, &object)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(object, "查询成功")
}
