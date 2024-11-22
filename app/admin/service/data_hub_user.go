package service

import (
	"fmt"
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
	orm := e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	tempSize := c.GetPageSize()
	offset := (c.GetPageIndex() - 1) * tempSize
	whereCondition := ""
	if c.Email != "" || c.TwitterName != "" || c.TelegramName != "" || c.SolanaAccount != "" || c.ID != 0 || c.UserName != "" {
		whereCondition += "WHERE "
		if c.Email != "" {
			whereCondition += "email= '" + fmt.Sprintf("%s", c.Email) + "' AND "
		}
		if c.TwitterName != "" {
			whereCondition += "twitter_name= '" + fmt.Sprintf("%s", c.TwitterName) + "' AND "
		}
		if c.TelegramName != "" {
			whereCondition += "telegram_name='" + fmt.Sprintf("%s", c.TelegramName) + "' AND "
		}
		if c.SolanaAccount != "" {
			whereCondition += "wallet= '" + fmt.Sprintf("%s", c.SolanaAccount) + "' AND "
		}
		if c.ID != 0 {
			whereCondition += "u.id= " + fmt.Sprintf("%d", c.ID) + " AND "
		}
		if c.UserName != "" {
			whereCondition += "name= '" + fmt.Sprintf("%s", c.UserName) + "' AND "
		}
		whereCondition += "1=1"
	}
	// 按字段排序
	orderCondition := "created_at ASC"
	if c.RankOrder != "" {
		if c.RankOrder == "descend" {
			orderCondition = "point_total ASC"
		} else {
			orderCondition = "point_total DESC"
		}
	} else if c.PointOrder != "" {
		if c.PointOrder == "descend" {
			orderCondition = "point_total DESC"
		} else {
			orderCondition = "point_total ASC"
		}
	} else if c.CompletedOrder != "" {
		if c.CompletedOrder == "descend" {
			orderCondition = "completed_items DESC"
		} else {
			orderCondition = "completed_items ASC"
		}
	}

	err = orm.Raw(fmt.Sprintf(`
SELECT id,rank, email, name,     point_total as point, wallet, twitter_name, level, location, telegram_name, telegram_full_name, near_account,evm_account,completed_items,upload_times,contribution_value,created_at
FROM (
    SELECT 
    ROW_NUMBER() OVER (ORDER BY point_total DESC) as rank,
    subquery.*
FROM (
    SELECT 
		    u.id, email, name, wallet, twitter_name, level, location, telegram_name, u.telegram_full_name, (SELECT near_address FROM user_near_addresses WHERE "user" = u.id) as near_account,
				(SELECT ethereum_address FROM user_ethereum_addresses WHERE "user" = u.id) as evm_account, (SELECT COUNT(*) FROM trains WHERE "user" = u.id) as completed_items,
				(SELECT COUNT(*) FROM ai_task_upload_records WHERE "user" = u.id) as upload_times,(SELECT COUNT(*) FROM ai_task_uploaded_files WHERE "user" = u.id and v_pass = true and a_pass = true) as contribution_value,
        (u.point + COALESCE(t.tma_point, 0)) as point,u.created_at,
        CASE 
            WHEN u.level = 1 THEN u.point + COALESCE(t.tma_point, 0)
            WHEN u.level = 2 THEN u.point + COALESCE(t.tma_point, 0) + %d
            ELSE u.point + COALESCE(t.tma_point, 0) + %d
        END AS point_total
    FROM users u
    LEFT JOIN tma_users t ON u.telegram_id = t.telegram_id %s
) as subquery ORDER BY %s
) AS subquery2 LIMIT ? OFFSET ?;
`, 40000, 120000+40000, whereCondition, orderCondition),
		tempSize, offset).Scan(list).Error

	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	err = orm.Raw(fmt.Sprintf(`
			SELECT COUNT(*)
			FROM (
				SELECT *
				FROM users u
				LEFT JOIN tma_users t ON u.telegram_id = t.telegram_id %s
			) as subquery;
		`, whereCondition)).Scan(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	//err = e.Orm.Debug().
	//	Scopes(
	//		cDto.MakeCondition(c.GetNeedSearch()),
	//		cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
	//		actions.Permission(data.TableName(), p),
	//	).Order("id asc").
	//	Find(list).Limit(-1).Offset(-1).
	//	Count(count).Error
	//if err != nil {
	//	e.Log.Errorf("db error: %s", err)
	//	return err
	//}
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

func (e *DataHubUser) GetPageAllPoint(c *dto.DataHubUserGetAllRewardReq, p *actions.DataPermission, list *[]models.AllRewardItem, count *int64, average *int64) error {
	var err error
	var data models.Train
	orm := e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	// 按字段排序
	orderCondition := "\"user\" ASC"
	if c.PointOrder != "" {
		if c.PointOrder == "descend" {
			orderCondition = "point DESC"
		} else {
			orderCondition = "point ASC"
		}
	}
	tempSize := c.GetPageSize()
	offset := (c.GetPageIndex() - 1) * tempSize
	err = orm.Raw(fmt.Sprintf("SELECT \"user\", COALESCE(SUM(point), 0) AS point FROM \"train_rewards\" WHERE point!=0 AND train_rewards.created_at >= to_timestamp(?) "+
		"AND train_rewards.created_at <= to_timestamp(?) GROUP BY \"user\" ORDER BY %s LIMIT ? OFFSET ?", orderCondition),
		c.StartTime, c.EndTime, tempSize, offset).Scan(list).Error

	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}

	err = orm.Raw("SELECT COUNT(*) FROM (SELECT \"user\" FROM \"train_rewards\" WHERE point!=0 AND train_rewards.created_at >= to_timestamp(?) "+
		"AND train_rewards.created_at <= to_timestamp(?) GROUP BY \"user\" ) AS subquery", c.StartTime, c.EndTime).Scan(count).Error

	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	err = orm.Raw("SELECT COALESCE(ROUND(AVG(point),0), 0) FROM (SELECT \"user\", COALESCE(SUM(point), 0) as point FROM \"train_rewards\" WHERE point!=0 AND train_rewards.created_at >= to_timestamp(?) "+
		"AND train_rewards.created_at <= to_timestamp(?) GROUP BY \"user\" ) AS subquery", c.StartTime, c.EndTime).Scan(average).Error

	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

func (e *DataHubUser) GetPageAllPointExport(c *dto.DataHubUserGetAllRewardReq, p *actions.DataPermission, list *[]models.AllRewardItem) error {
	var err error
	var data models.Train
	orm := e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	// 按字段排序
	orderCondition := "\"user\" ASC"
	if c.PointOrder != "" {
		if c.PointOrder == "descend" {
			orderCondition = "point DESC"
		} else {
			orderCondition = "point ASC"
		}
	}

	err = orm.Raw(fmt.Sprintf("SELECT \"user\", COALESCE(SUM(point), 0) AS point FROM \"train_rewards\" WHERE point!=0 AND train_rewards.created_at >= to_timestamp(?) "+
		"AND train_rewards.created_at <= to_timestamp(?) GROUP BY \"user\" ORDER BY %s", orderCondition),
		c.StartTime, c.EndTime).Scan(list).Error

	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}

	return nil
}
