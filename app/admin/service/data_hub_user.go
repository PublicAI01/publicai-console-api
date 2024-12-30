package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/config"
	"io/ioutil"
	"net/http"
	"net/url"
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
			whereCondition += "email ILIKE '" + fmt.Sprintf("%s", c.Email) + "' AND "
		}
		if c.TwitterName != "" {
			whereCondition += "twitter_name ILIKE '" + fmt.Sprintf("%s", c.TwitterName) + "' AND "
		}
		if c.TelegramName != "" {
			whereCondition += "telegram_name ILIKE '" + fmt.Sprintf("%s", c.TelegramName) + "' AND "
		}
		if c.SolanaAccount != "" {
			whereCondition += "wallet= '" + fmt.Sprintf("%s", c.SolanaAccount) + "' AND "
		}
		if c.ID != 0 {
			whereCondition += "u.id= " + fmt.Sprintf("%d", c.ID) + " AND "
		}
		if c.UserName != "" {
			whereCondition += "name ILIKE '" + fmt.Sprintf("%s", c.UserName) + "' AND "
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
SELECT id,rank, email, name,     point_total as point, wallet, twitter_name, level, location, telegram_name, ambassador, telegram_full_name, near_account,evm_account,completed_items,upload_times,contribution_value,created_at
FROM (
    SELECT 
    ROW_NUMBER() OVER (ORDER BY point_total DESC) as rank,
    subquery.*
FROM (
    SELECT 
		    u.id, email, name, wallet, twitter_name, level, location, telegram_name, ambassador, u.telegram_full_name, (SELECT near_address FROM user_near_addresses WHERE "user" = u.id) as near_account,
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
	orderCondition := "created_at ASC"
	if c.UID == 0 {
		//err = orm.Raw("SELECT id, \"user\",point, created_at FROM \"train_rewards\" WHERE train_rewards.created_at >= ? "+
		//	"AND train_rewards.created_at <= ? LIMIT ? OFFSET ?",
		//	c.StartTime, c.EndTime, tempSize, offset).Scan(list).Error
		err = orm.Raw(fmt.Sprintf(`
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ? ORDER BY %s LIMIT ? OFFSET ?
`, orderCondition),
			c.StartTime, c.EndTime,
			c.StartTime, c.EndTime,
			c.StartTime, c.EndTime,
			tempSize, offset).Scan(list).Error
	} else {
		//err = orm.Raw("SELECT id, \"user\",point, created_at FROM \"train_rewards\" WHERE train_rewards.user = ? AND train_rewards.created_at >= ? "+
		//	"AND train_rewards.created_at <= ? LIMIT ? OFFSET ?",
		//	c.UID, c.StartTime, c.EndTime, tempSize, offset).Scan(list).Error
		err = orm.Raw(fmt.Sprintf(`
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ? AND "user"=?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ?  AND "user"=?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id AND users."id"=?) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ? ORDER BY %s LIMIT ? OFFSET ?
`, orderCondition),
			c.StartTime, c.EndTime, c.UID,
			c.StartTime, c.EndTime, c.UID,
			c.UID, c.StartTime, c.EndTime,
			tempSize, offset).Scan(list).Error
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if c.UID == 0 {
		//err = orm.Raw("SELECT COUNT(*) FROM \"train_rewards\" WHERE train_rewards.created_at >= ? "+
		//	"AND train_rewards.created_at <= ?", c.StartTime, c.EndTime).Scan(count).Error
		err = orm.Raw(`
SELECT 
    COUNT(*) AS total_rows 
FROM (
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ?) AS grouped_users;
`,
			c.StartTime, c.EndTime,
			c.StartTime, c.EndTime,
			c.StartTime, c.EndTime).Scan(count).Error
	} else {
		//err = orm.Raw("SELECT COUNT(*) FROM \"train_rewards\" WHERE train_rewards.user = ? AND train_rewards.created_at >= ? "+
		//	"AND train_rewards.created_at <= ?", c.UID, c.StartTime, c.EndTime).Scan(count).Error
		err = orm.Raw(`
SELECT 
    COUNT(*) AS total_rows 
FROM (
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ? AND "user"=?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ? AND "user"=?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id AND users."id"=?) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ?) AS grouped_users;
`,
			c.StartTime, c.EndTime, c.UID,
			c.StartTime, c.EndTime, c.UID,
			c.UID, c.StartTime, c.EndTime).Scan(count).Error
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if c.UID == 0 {
		//err = orm.Raw("SELECT COALESCE(SUM(point), 0) AS point FROM \"train_rewards\" WHERE train_rewards.created_at >= ? "+
		//	"AND train_rewards.created_at <= ?", c.StartTime, c.EndTime).Scan(total).Error
		err = orm.Raw(`SELECT COALESCE(SUM(point), 0) FROM (SELECT "user",COALESCE(SUM(sub_query.point), 0) AS point  FROM (
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ?) as sub_query   GROUP BY "user") AS sub_query`,
			c.StartTime, c.EndTime,
			c.StartTime, c.EndTime,
			c.StartTime, c.EndTime).Scan(total).Error
	} else {
		//err = orm.Raw("SELECT COALESCE(SUM(point), 0) AS point FROM \"train_rewards\" WHERE train_rewards.user = ? AND train_rewards.created_at >= ? "+
		//	"AND train_rewards.created_at <= ?", c.UID, c.StartTime, c.EndTime).Scan(total).Error
		err = orm.Raw(`SELECT COALESCE(SUM(point), 0) FROM (SELECT "user",COALESCE(SUM(sub_query.point), 0) AS point  FROM (
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ? AND "user"=?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ? AND "user"=?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id AND users."id"=?) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ?) as sub_query   GROUP BY "user") AS sub_query`,
			c.StartTime, c.EndTime, c.UID,
			c.StartTime, c.EndTime, c.UID,
			c.UID, c.StartTime, c.EndTime).Scan(total).Error
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

func (e *DataHubUser) GetPageUserPointExport(c *dto.DataHubUserPointGetPageReq, p *actions.DataPermission, list *[]models.RewardItem, count *int64, total *int64) error {
	var err error
	var data models.Train
	orm := e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	orderCondition := "created_at ASC"
	if c.UID == 0 {
		err = orm.Raw(fmt.Sprintf(`
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ? ORDER BY %s
`, orderCondition),
			c.StartTime, c.EndTime,
			c.StartTime, c.EndTime,
			c.StartTime, c.EndTime,
		).Scan(list).Error
	} else {
		err = orm.Raw(fmt.Sprintf(`
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ? AND "user"=?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ?  AND "user"=?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id AND users."id"=?) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ? ORDER BY %s
`, orderCondition),
			c.StartTime, c.EndTime, c.UID,
			c.StartTime, c.EndTime, c.UID,
			c.UID, c.StartTime, c.EndTime,
		).Scan(list).Error
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if c.UID == 0 {
		err = orm.Raw(`
SELECT 
    COUNT(*) AS total_rows 
FROM (
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ?) AS grouped_users;
`,
			c.StartTime, c.EndTime,
			c.StartTime, c.EndTime,
			c.StartTime, c.EndTime).Scan(count).Error
	} else {
		err = orm.Raw(`
SELECT 
    COUNT(*) AS total_rows 
FROM (
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ? AND "user"=?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ? AND "user"=?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id AND users."id"=?) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ?) AS grouped_users;
`,
			c.StartTime, c.EndTime, c.UID,
			c.StartTime, c.EndTime, c.UID,
			c.UID, c.StartTime, c.EndTime).Scan(count).Error
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if c.UID == 0 {
		err = orm.Raw(`SELECT COALESCE(SUM(point), 0) FROM (SELECT "user",COALESCE(SUM(sub_query.point), 0) AS point  FROM (
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ?) as sub_query   GROUP BY "user") AS sub_query`,
			c.StartTime, c.EndTime,
			c.StartTime, c.EndTime,
			c.StartTime, c.EndTime).Scan(total).Error
	} else {
		err = orm.Raw(`SELECT COALESCE(SUM(point), 0) FROM (SELECT "user",COALESCE(SUM(sub_query.point), 0) AS point  FROM (
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ? AND "user"=?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ? AND "user"=?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id AND users."id"=?) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ?) as sub_query   GROUP BY "user") AS sub_query`,
			c.StartTime, c.EndTime, c.UID,
			c.StartTime, c.EndTime, c.UID,
			c.UID, c.StartTime, c.EndTime).Scan(total).Error
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
	err = orm.Raw(fmt.Sprintf(`
SELECT "user",COALESCE(SUM(sub_query.point), 0) AS point  FROM (
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ?) as sub_query   GROUP BY "user" ORDER BY %s LIMIT ? OFFSET ?
`, orderCondition),
		c.StartTime, c.EndTime,
		c.StartTime, c.EndTime,
		c.StartTime, c.EndTime,
		tempSize, offset).Scan(list).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}

	err = orm.Raw(`
SELECT 
    COUNT(*) AS total_rows 
FROM (SELECT "user",COALESCE(SUM(sub_query.point), 0) AS point  FROM (
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ?) as sub_query   GROUP BY "user") AS grouped_users;
`,
		c.StartTime, c.EndTime,
		c.StartTime, c.EndTime,
		c.StartTime, c.EndTime).Scan(count).Error

	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	err = orm.Raw(`SELECT COALESCE(ROUND(AVG(point),0), 0) FROM (SELECT "user",COALESCE(SUM(sub_query.point), 0) AS point  FROM (
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ?) as sub_query   GROUP BY "user") AS sub_query`,
		c.StartTime, c.EndTime,
		c.StartTime, c.EndTime,
		c.StartTime, c.EndTime).Scan(average).Error

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

	err = orm.Raw(fmt.Sprintf(`
SELECT "user",COALESCE(SUM(sub_query.point), 0) AS point  FROM (
SELECT "user",  point, created_at FROM train_rewards WHERE point != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT "user", point_reward as point, created_at FROM ai_task_rewards WHERE point_reward != 0 AND created_at >= ? AND created_at <= ?
UNION ALL SELECT B."user",A.point, A.created_at FROM tma_rewards as A INNER JOIN (SELECT tma_users."id" as "id",users."id" as "user" FROM tma_users  
INNER JOIN users ON tma_users.telegram_id=users.telegram_id) as B ON A.tma_id=B."id"  
WHERE A.point != 0 AND created_at >= ? AND created_at <= ?) as sub_query   GROUP BY "user" ORDER BY %s
`, orderCondition),
		c.StartTime, c.EndTime,
		c.StartTime, c.EndTime,
		c.StartTime, c.EndTime).Scan(list).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}

	return nil
}

func (e *DataHubUser) GetPageAmbassador(c *dto.DataHubAmbassadorGetPageReq, p *actions.DataPermission, list *[]models.DataHubAmbassador, count *int64) error {
	var err error
	var data models.DataHubAmbassador
	orm := e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	tempSize := c.GetPageSize()
	offset := (c.GetPageIndex() - 1) * tempSize
	whereCondition := "WHERE "
	if c.ID != 0 {
		if c.ID != 0 {
			whereCondition += "u.id= " + fmt.Sprintf("%d", c.ID) + " AND "
		}

	}
	whereCondition += " ambassador=TRUE "
	// 按字段排序
	orderCondition := "consensus_contribution DESC"
	if c.ContributionOrder != "" {
		if c.ContributionOrder == "descend" {
			orderCondition = "consensus_contribution DESC"
		} else {
			orderCondition = "consensus_contribution ASC"
		}
	}

	err = orm.Raw(fmt.Sprintf(`SELECT id, email, name,  twitter_name, location, telegram_name, telegram_full_name,
	      (SELECT COUNT(*) FROM ai_task_trains WHERE "user" = u.id) as consensus_contribution FROM users u
	       %s ORDER BY %s LIMIT ? OFFSET ?;
	`, whereCondition, orderCondition),
		tempSize, offset).Scan(list).Error
	//	err = orm.Raw(fmt.Sprintf(`SELECT
	//    u.id,
	//    u.email,
	//    u.name,
	//    u.twitter_name,
	//    u.location,
	//    u.telegram_name,
	//    u.telegram_full_name,
	//    COUNT(a.id) AS consensus_contribution
	//FROM
	//    users u %s
	//LEFT JOIN
	//    ai_task_trains a ON a."user" = u.id
	//GROUP BY
	//    u.id, u.email, u.name, u.twitter_name, u.location, u.telegram_name, u.telegram_full_name
	//ORDER BY %s LIMIT ? OFFSET ?`, whereCondition, orderCondition), tempSize, offset).Scan(list).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	err = orm.Raw(fmt.Sprintf(`
			SELECT COUNT(*)
			FROM users %s;
		`, whereCondition)).Scan(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

func (e *DataHubUser) UpdateAmbassador(c *dto.MarketplaceAmbassadorUpdateReq, p *actions.DataPermission) error {
	hubUrl := fmt.Sprintf("%s/api/admin/user/ambassador", config.ExtConfig.DataHubIp)
	params := url.Values{}
	params.Set("token", config.ExtConfig.Token)
	urlWithParams := fmt.Sprintf("%s?%s", hubUrl, params.Encode())
	postBody, _ := json.Marshal(map[string]interface{}{
		"ambassadors": c.Ambassadors,
	})
	// 将数据转换为字节序列
	requestBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest("PUT", urlWithParams, requestBody)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	defer response.Body.Close()
	type PutResponse struct {
		Data interface{} `json:"data"`
		Msg  string      `json:"msg"`
		Code int         `json:"code"`
	}
	var putResp PutResponse
	if err == nil {
		body, readErr := ioutil.ReadAll(response.Body)
		fmt.Println(body)
		if readErr == nil {
			err = json.Unmarshal(body, &putResp)
			if err == nil && putResp.Code == 200 {
				return nil
			}
		} else {
			err = readErr
		}
	}

	return err
}

func (e *DataHubUser) GetPageAmbassadorExport(c *dto.DataHubExportAmbassadorReq, p *actions.DataPermission, list *[]models.DataHubAmbassador) error {
	var err error
	var data models.Train
	orm := e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	// 按字段排序
	orderCondition := "consensus_contribution DESC"
	if c.ContributionOrder != "" {
		if c.ContributionOrder == "descend" {
			orderCondition = "consensus_contribution DESC"
		} else {
			orderCondition = "consensus_contribution ASC"
		}
	}

	err = orm.Raw(fmt.Sprintf(`SELECT id, email, name,  twitter_name, location, telegram_name, telegram_full_name, 
       (SELECT COUNT(*) FROM ai_task_trains WHERE "user" = u.id) as consensus_contribution FROM users as u 
        where  ambassador=TRUE ORDER BY %s;
`, orderCondition)).Scan(list).Error

	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}

	return nil
}
