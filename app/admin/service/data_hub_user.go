package service

import (
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type DataHubUser struct {
	service.Service
}

func (e *DataHubUser) GetPageUser(c *dto.DataHubUserGetPageReq, p *actions.DataPermission, list *[]models.DataHubUser, count *int64) error {
	var err error
	var data models.DataHubUser
	err = e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order("id asc").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

func (e *DataHubUser) GetPageUserPoint(c *dto.DataHubUserPointGetPageReq, p *actions.DataPermission, list *[]models.RewardItem, count *int64, total *int64) error {
	var err error
	var data models.Train
	orm := e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	tempSize := c.GetPageSize()
	offset := (c.GetPageIndex() - 1) * tempSize
	if c.UID == 0 {
		err = orm.Raw("SELECT id, \"user\",point, created_at FROM \"train_rewards\" WHERE train_rewards.created_at >= ? "+
			"AND train_rewards.created_at <= ? LIMIT ? OFFSET ?",
			c.StartTime, c.EndTime, tempSize, offset).Scan(list).Error
	} else {
		err = orm.Raw("SELECT id, \"user\",point, created_at FROM \"train_rewards\" WHERE train_rewards.user = ? AND train_rewards.created_at >= ? "+
			"AND train_rewards.created_at <= ? LIMIT ? OFFSET ?",
			c.UID, c.StartTime, c.EndTime, tempSize, offset).Scan(list).Error
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if c.UID == 0 {
		err = orm.Raw("SELECT COUNT(*) FROM \"train_rewards\" WHERE train_rewards.created_at >= ? "+
			"AND train_rewards.created_at <= ?", c.StartTime, c.EndTime).Scan(count).Error
	} else {
		err = orm.Raw("SELECT COUNT(*) FROM \"train_rewards\" WHERE train_rewards.user = ? AND train_rewards.created_at >= ? "+
			"AND train_rewards.created_at <= ?", c.UID, c.StartTime, c.EndTime).Scan(count).Error
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if c.UID == 0 {
		err = orm.Raw("SELECT COALESCE(SUM(point), 0) AS point FROM \"train_rewards\" WHERE train_rewards.created_at >= ? "+
			"AND train_rewards.created_at <= ?", c.StartTime, c.EndTime).Scan(total).Error
	} else {
		err = orm.Raw("SELECT COALESCE(SUM(point), 0) AS point FROM \"train_rewards\" WHERE train_rewards.user = ? AND train_rewards.created_at >= ? "+
			"AND train_rewards.created_at <= ?", c.UID, c.StartTime, c.EndTime).Scan(total).Error
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}
