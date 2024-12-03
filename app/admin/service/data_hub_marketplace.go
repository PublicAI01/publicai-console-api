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
		).Order("created_at desc").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

func (e *DataHubMarketplace) GetCampaignValidation(c *dto.DataHubMarketplaceGetCampaignValidationReq, p *actions.DataPermission, list *[]models.AITaskShowRecordItem, count *int64) error {
	var err error
	var data models.AITaskUploadRecord
	orm := e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	fmt.Println(p.UserId)
	tempSize := c.GetPageSize()
	offset := (c.GetPageIndex() - 1) * tempSize
	// 按字段排序
	orderCondition := "upload_time ASC"
	if c.CreatedAtOrder != "" {
		if c.CreatedAtOrder == "descend" {
			orderCondition = "upload_time DESC"
		} else {
			orderCondition = "upload_time ASC"
		}
	}
	err = orm.Raw(fmt.Sprintf(`
SELECT 
    ROW_NUMBER() OVER() AS no,
    u."id",
    u.success as data_number,
    u."user",
    u.status,
    COALESCE((SELECT CASE 
            WHEN ai_operate = TRUE THEN 1
            WHEN ai_operate = FALSE THEN 0
    END
    FROM ai_task_validations WHERE upload_record = u."id"
                ORDER BY created_at DESC
    LIMIT 1 )  ,-1) as editor,
    (SELECT COUNT(*) 
     FROM ai_task_uploaded_files 
     WHERE upload_record = u."id" AND a_pass = TRUE) AS valid,
    u.created_at AS upload_time
FROM 
    ai_task_upload_records AS u 
WHERE u.task= ? ORDER BY %s LIMIT ? OFFSET ?`, orderCondition), c.TaskID, tempSize, offset).Scan(&list).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	err = orm.Raw("SELECT COUNT(*) FROM ai_task_upload_records WHERE task = ?  ", c.TaskID).Scan(count).Error
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
	hubUrl := fmt.Sprintf("%s/api/admin/marketplace/campaign/validation", config.ExtConfig.DataHubIp)
	params := url.Values{}
	params.Set("token", config.ExtConfig.Token)
	urlWithParams := fmt.Sprintf("%s?%s", hubUrl, params.Encode())
	postBody, _ := json.Marshal(map[string]interface{}{
		"validations": c.Validations,
		"uid":         p.UserId,
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

func (e *DataHubMarketplace) GetCampaignDispute(c *dto.DataHubMarketplaceGetCampaignDisputeReq, p *actions.DataPermission, list *[]models.AITaskShowDisputeRecordItem, count *int64) error {
	var err error
	var data models.AITaskUploadRecord
	orm := e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	tempSize := c.GetPageSize()
	offset := (c.GetPageIndex() - 1) * tempSize
	// TODO:check completed
	err = orm.Raw(`
SELECT 
    ROW_NUMBER() OVER() AS no,
    u."id",
    u.success as data_number,
    u."user",
    u.status,
    COALESCE((SELECT CASE 
            WHEN ai_operate = TRUE THEN 1
            WHEN ai_operate = FALSE THEN 0
    END
    FROM ai_task_validations WHERE upload_record = u."id"
                ORDER BY created_at DESC
    LIMIT 1 )  ,-1) as editor,
    (SELECT COUNT(*) 
     FROM ai_task_uploaded_files 
     WHERE upload_record = u."id" AND a_pass = TRUE) AS valid,
    u.created_at AS upload_time
FROM 
    ai_task_upload_records AS u 
WHERE u.task= ? and completed=2 and can_issue = false LIMIT ? OFFSET ?`, c.TaskID, tempSize, offset).Scan(&list).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	err = orm.Raw("SELECT COUNT(*) FROM ai_task_upload_records WHERE task = ?  and completed=2 and can_issue = false  ", c.TaskID).Scan(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

func (e *DataHubMarketplace) UpdateDispute(c *dto.MarketplaceDisputeUpdateReq, p *actions.DataPermission) error {
	hubUrl := fmt.Sprintf("%s/api/admin/marketplace/campaign/dispute", config.ExtConfig.DataHubIp)
	params := url.Values{}
	params.Set("token", config.ExtConfig.Token)
	urlWithParams := fmt.Sprintf("%s?%s", hubUrl, params.Encode())
	postBody, _ := json.Marshal(map[string]interface{}{
		"task_id":     c.TaskID,
		"validations": c.Validations,
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
