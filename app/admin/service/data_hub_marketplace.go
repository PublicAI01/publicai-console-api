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

func (e *DataHubMarketplace) GetPageCampaignVariants(c *dto.DataHubMarketplaceGetPageCampaignReq, p *actions.DataPermission, list *[]models.AITaskVariants, count *int64) error {
	var err error
	var data models.AITask
	orm := e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	tempSize := c.GetPageSize()
	offset := (c.GetPageIndex() - 1) * tempSize

	err = orm.Raw(fmt.Sprintf(`SELECT 
    t."id",
    t."name",
    t."start",
    t."end",
    COALESCE(upload_records.upload_times, 0) AS upload_times,
    COALESCE(upload_records.total_upload_users, 0) AS total_upload_users,
    COALESCE(upload_files.amount_of_all, 0) AS amount_of_all,
    COALESCE(accepted_records.accepted_times, 0) AS accepted_times,
    COALESCE(upload_files.accepted_total, 0) AS accepted_total,
    COALESCE(accepted_records.accepted_users, 0) AS accepted_users
FROM ai_tasks AS t 
LEFT JOIN (
    SELECT 
        "task", 
        COUNT(*) AS upload_times,
        COUNT(DISTINCT "user") AS total_upload_users
    FROM 
        ai_task_upload_records 
    GROUP BY 
        "task"
) AS upload_records ON upload_records."task" = t."id"
LEFT JOIN (
    SELECT 
        "task", 
        COUNT(*) AS amount_of_all,
				COUNT(CASE WHEN status = 1 THEN 1 END) AS accepted_total
    FROM 
        ai_task_uploaded_files 
    GROUP BY 
        "task"
) AS upload_files ON upload_files."task" = t."id"
LEFT JOIN (
    SELECT 
        "task",
        COUNT(*) AS accepted_times,
        COUNT(DISTINCT "user") AS accepted_users
    FROM 
        ai_task_upload_records 
    WHERE 
        status > 1
    GROUP BY 
        "task"
) AS accepted_records ON accepted_records."task" = t."id" ORDER BY t."created_at" DESC LIMIT ? OFFSET ?`), tempSize, offset).Scan(&list).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	err = orm.Raw("SELECT COUNT(*) FROM ai_tasks").Scan(count).Error
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
	timeCondition := ""
	if c.StartTime != "" && c.EndTime != "" {
		timeCondition = fmt.Sprintf("and u.created_at >='%s' and u.created_at <= '%s'", c.StartTime, c.EndTime)
	}
	userCondition := ""
	if c.UID != 0 {
		userCondition = fmt.Sprintf("and \"user\"=%d", c.UID)
	}
	statusCondition := ""
	if c.UID != 0 {
		statusCondition = fmt.Sprintf("and status=%d", c.Status)
	}
	type Consensus struct {
		Consensus int `json:"consensus"`
	}
	var consensus Consensus

	err = orm.Raw(`SELECT consensus from ai_tasks where id =?`, c.TaskID).Scan(&consensus).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	findConsensus := consensus.Consensus/2 + 1
	err = orm.Raw(fmt.Sprintf(`
SELECT 
    ROW_NUMBER() OVER() AS no,
    u."id",
    u.total as data_number,
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
     WHERE upload_record = u."id" AND (v_aye >= %d OR v_nay >=%d)) AS valid,
    u.created_at AS upload_time
FROM 
    ai_task_upload_records AS u 
WHERE u.task= ? %s %s %s ORDER BY %s LIMIT ? OFFSET ?`, findConsensus, findConsensus, userCondition, timeCondition, statusCondition, orderCondition), c.TaskID, tempSize, offset).Scan(&list).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	err = orm.Raw(fmt.Sprintf("SELECT COUNT(*) FROM ai_task_upload_records as u WHERE u.task = ?  %s %s %s", timeCondition, userCondition, statusCondition), c.TaskID).Scan(count).Error
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
func (e *DataHubMarketplace) UpdateMalicious(c *dto.MarketplaceValidationMaliciousUpdateReq, p *actions.DataPermission) error {
	hubUrl := fmt.Sprintf("%s/api/admin/marketplace/campaign/validation/malicious", config.ExtConfig.DataHubIp)
	params := url.Values{}
	params.Set("token", config.ExtConfig.Token)
	urlWithParams := fmt.Sprintf("%s?%s", hubUrl, params.Encode())
	postBody, _ := json.Marshal(map[string]interface{}{
		"id":     c.Id,
		"m_flag": c.Flag,
		"uid":    p.UserId,
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
	type Consensus struct {
		Consensus int `json:"consensus"`
	}
	var consensus Consensus

	err = orm.Raw(`SELECT consensus from ai_tasks where id =?`, c.TaskID).Scan(&consensus).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	userCondition := ""
	if c.UID != 0 {
		userCondition = fmt.Sprintf("and \"user\"=%d", c.UID)
	}
	findConsensus := consensus.Consensus/2 + 1
	needAm := consensus.Consensus / 2 / 5
	// TODO:check completed
	err = orm.Raw(fmt.Sprintf(`
SELECT 
    ROW_NUMBER() OVER() AS no,
    u."id",
    u.total as data_number,
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
     WHERE upload_record = u."id" AND (v_aye >= %d OR v_nay >=%d)) AS valid,
    u.created_at AS upload_time
FROM 
    ai_task_upload_records AS u 
WHERE u.task= ? and completed=? and can_issue = false %s LIMIT ? OFFSET ?`, findConsensus, findConsensus, userCondition), c.TaskID, needAm, tempSize, offset).Scan(&list).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	err = orm.Raw(fmt.Sprintf("SELECT COUNT(*) FROM ai_task_upload_records WHERE task = ?  and completed=? and can_issue = false  %s", userCondition), c.TaskID, needAm).Scan(count).Error
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

func (e *DataHubMarketplace) GetValidationSummary(c *dto.GetCampaignValidationSummaryReq, p *actions.DataPermission, model *models.ValidationSummary) *DataHubMarketplace {
	var data models.ValidationSummary
	orm := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		)
	userCondition := ""
	if c.UID != 0 {
		userCondition = fmt.Sprintf("and \"user\"=%d", c.UID)
	}
	statusCondition := ""
	if c.UID != 0 {
		statusCondition = fmt.Sprintf("and status=%d", c.Status)
	}
	err := orm.Raw(fmt.Sprintf(`
SELECT COUNT(*) as total ,
SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) as all_accepted,
SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) as partial_accepted,
SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END) as pending,
SUM(CASE WHEN status = -1 THEN 1 ELSE 0 END) as subpar,
SUM(CASE WHEN status = -2 THEN 1 ELSE 0 END) as malicious
FROM ai_task_upload_records WHERE task=? AND created_at >= ? AND created_at <= ? %s %s;
`, userCondition, statusCondition), c.TaskID, c.StartTime, c.EndTime).Scan(&model).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		_ = e.AddError(err)
		return e
	}
	return e
}

func (e *DataHubMarketplace) DownloadValidation(c *dto.GetCampaignValidationSummaryReq, p *actions.DataPermission, list *[]models.AITaskShowRecordItem) error {
	var err error
	var data models.AITaskUploadRecord
	orm := e.Orm.Debug().
		Scopes(
			actions.Permission(data.TableName(), p),
		)
	// 按字段排序
	orderCondition := "upload_time ASC"
	timeCondition := ""
	if c.StartTime != "" && c.EndTime != "" {
		timeCondition = fmt.Sprintf("and u.created_at >='%s' and u.created_at <= '%s'", c.StartTime, c.EndTime)
	}
	userCondition := ""
	if c.UID != 0 {
		userCondition = fmt.Sprintf("and \"user\"=%d", c.UID)
	}
	statusCondition := ""
	if c.UID != 0 {
		statusCondition = fmt.Sprintf("and status=%d", c.Status)
	}
	type Consensus struct {
		Consensus int `json:"consensus"`
	}
	var consensus Consensus

	err = orm.Raw(`SELECT consensus from ai_tasks where id =?`, c.TaskID).Scan(&consensus).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	findConsensus := consensus.Consensus/2 + 1
	err = orm.Raw(fmt.Sprintf(`
SELECT 
    ROW_NUMBER() OVER() AS no,
    u."id",
    u.total as data_number,
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
     WHERE upload_record = u."id" AND (v_aye >= %d OR v_nay >=%d)) AS valid,
    u.created_at AS upload_time
FROM 
    ai_task_upload_records AS u 
WHERE u.task= ? %s %s %s ORDER BY %s`, findConsensus, findConsensus, userCondition, timeCondition, statusCondition, orderCondition), c.TaskID).Scan(&list).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}
