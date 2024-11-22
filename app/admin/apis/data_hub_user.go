package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/xuri/excelize/v2"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"strconv"
	"time"
)

type DataHubUser struct {
	api.Api
}

// GetPage
// @Summary 列表DataHub用户信息数据
// @Description 获取JSON
// @Tags DataHub
// @Param email query string false "email"
// @Param telegram_name query string false "telegram_name"
// @Param twitter_name query string false "twitter_name"
// @Param solana_account query string false "solana_account"
// @Success 200 {string} {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/data_hub/user [get]
// @Security Bearer
func (e DataHubUser) GetPageUser(c *gin.Context) {
	s := service.DataHubUser{}
	req := dto.DataHubUserGetPageReq{}
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

	list := make([]models.DataHubUser, 0)
	var count int64

	err = s.GetPageUser(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	for i, user := range list {
		var tgName *string
		if user.TelegramName == nil {
			tgName = user.TelegramFullName
		} else {
			tgName = user.TelegramName
		}
		list[i].TelegramName = tgName
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetPage
// @Summary 列表DataHub用户积分信息数据
// @Description 获取JSON
// @Tags DataHub
// @Param uid query string false "uid"
// @Param start_time query string false "start_time"
// @Param end_time query string false "end_time"
// @Success 200 {string} {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/data_hub/user/point [get]
// @Security Bearer
func (e DataHubUser) GetPageUserPoint(c *gin.Context) {
	s := service.DataHubUser{}
	req := dto.DataHubUserPointGetPageReq{}
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

	list := make([]models.RewardItem, 0)
	var count int64
	var total int64

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
	err = s.GetPageUserPoint(&req, p, &list, &count, &total)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(map[string]interface{}{
		"total": total,
		"data":  list,
	}, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetAllReward Get
// @Summary 获取user reward信息数据
// @Description 获取JSON
// @Tags DataHub
// @Param start_time query string false "start_time"
// @Param end_time query string false "end_time"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/data_hub/user/reward [get]
// @Security Bearer
func (e DataHubUser) GetAllReward(c *gin.Context) {
	s := service.DataHubUser{}
	req := dto.DataHubUserGetAllRewardReq{}
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

	list := make([]models.AllRewardItem, 0)
	var count int64
	var average int64

	err = s.GetPageAllPoint(&req, p, &list, &count, &average)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}

	e.PageOK(map[string]interface{}{
		"average": average,
		"data":    list,
	}, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetAllRewardExport Get
// @Summary 导出user reward信息数据
// @Description 获取JSON
// @Tags DataHub
// @Param start_time query string false "start_time"
// @Param end_time query string false "end_time"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/data_hub/user/reward/export [get]
// @Security Bearer
func (e DataHubUser) GetAllRewardExport(c *gin.Context) {
	s := service.DataHubUser{}
	req := dto.DataHubUserGetAllRewardReq{}
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

	list := make([]models.AllRewardItem, 0)

	err = s.GetPageAllPointExport(&req, p, &list)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("points counting")
	if err != nil {
		e.Error(500, err, "导出失败")
		return
	}
	// Set value of a cell.
	f.SetCellValue("points counting", "A1", "User ID")
	f.SetCellValue("points counting", "B1", "Points")
	for i, item := range list {
		f.SetCellValue("points counting", fmt.Sprintf("A%d", i+2), item.User)
		f.SetCellValue("points counting", fmt.Sprintf("B%d", i+2), item.Point)
	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")
	//e.PageOK(map[string]interface{}{
	//	"average": average,
	//	"data":    list,
	//}, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	filename := fmt.Sprintf("points_%s_%s.xlsx", req.StartTime, req.EndTime)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	if err := f.Write(c.Writer); err != nil {
		e.Error(500, err, "导出失败")
		return
	}
	//e.OK(0, "export success")
}
