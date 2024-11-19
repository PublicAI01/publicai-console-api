package service

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type DataHubMarketplace struct {
	service.Service
}

func (e *DataHubMarketplace) GetPageCampaign(c *dto.DataHubMarketplaceGetPageCampaignReq, p *actions.DataPermission, list *[]models.AITask, count *int64) error {
	var err error
	var data models.AITask
	err = e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

func (e *DataHubMarketplace) GetCampaignValidation(c *dto.DataHubMarketplaceGetCampaignValidationReq, p *actions.DataPermission, list *[]models.AITaskUploadedFile, count *int64) error {
	var err error
	var data models.AITaskUploadedFile
	err = e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

func (e *DataHubMarketplace) GetCampaignReward(c *dto.DataHubMarketplaceGetPageRewardReq, p *actions.DataPermission, list *[]models.MarketPlaceRewardItem, count *int64) error {
	var err error
	var data models.AITaskRecord
	orm := e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	tempSize := c.GetPageSize()
	offset := (c.GetPageIndex() - 1) * tempSize
	err = orm.Raw("SELECT   ROW_NUMBER() OVER() AS no , \"user\", email, wallet as solana_account, usdt_reward, point_reward, "+
		"(SELECT SUM(success) FROM ai_task_upload_records WHERE \"user\" =ai_task_records.\"user\"  and task =?) as collected, "+
		"(SELECT COUNT(*) FROM ai_task_uploaded_files WHERE \"user\" =ai_task_records.\"user\"  and task =? and status = 1) as accepted, 0 as indicators, usdt_link, ai_task_records.created_at FROM ai_task_records INNER JOIN users ON ai_task_records.user = users.id WHERE task= ? LIMIT ? OFFSET ?",
		c.TaskID, c.TaskID, c.TaskID, tempSize, offset).Scan(list).Limit(-1).Offset(-1).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	err = orm.Raw("SELECT COUNT(*) FROM ai_task_records  INNER JOIN users ON ai_task_records.user = users.id WHERE task = ?  ", c.TaskID).Scan(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

func (e *DataHubMarketplace) Update(c *dto.MarketplaceValidationUpdateReq, p *actions.DataPermission) error {
	var model = models.AITaskUploadedFile{}
	db := e.Orm.Debug().First(&model, c.Id)
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	model.APass = c.APass
	status := -1
	if c.APass {
		status = 0
	}
	model.Status = status
	db = e.Orm.Save(&model)
	if err := db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysApi error:%s", err)
		return err
	}

	return nil
}
