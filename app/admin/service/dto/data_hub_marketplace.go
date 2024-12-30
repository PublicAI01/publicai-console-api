package dto

import (
	"go-admin/common/dto"
	"mime/multipart"
)

type DataHubMarketplaceGetPageCampaignReq struct {
	dto.Pagination `search:"-"`
	Email          string `form:"email" search:"type:contains;column:email;table:users"`
}

func (m *DataHubMarketplaceGetPageCampaignReq) GetNeedSearch() interface{} {
	return *m
}

type DataHubMarketplaceGetCampaignValidationReq struct {
	dto.Pagination `search:"-"`
	TaskID         int    `form:"task_id" search:"type:exact;column:task;table:ai_task_upload_records"`
	StartTime      string `form:"start_time" search:"type:gte;column:created_at;table:ai_task_upload_records"`
	EndTime        string `form:"end_time" search:"type:lte;column:created_at;table:ai_task_upload_records"`
	UID            int    `form:"uid" search:"type:exact;column:user;table:ai_task_upload_records"`
	Status         string `form:"status" search:"type:exact;column:status;table:ai_task_upload_records"`
	GetCampaignValidationOrder
}

type GetCampaignValidationOrder struct {
	CreatedAtOrder string `search:"type:order;column:created_at;table:ai_task_upload_records" form:"upload_timeOrder"`
}

func (m *DataHubMarketplaceGetCampaignValidationReq) GetNeedSearch() interface{} {
	return *m
}

type DataHubMarketplaceGetPageRewardReq struct {
	dto.Pagination `search:"-"`
	TaskID         int `form:"task_id" search:"type:exact;column:task;table:ai_task_records"`
}

func (m *DataHubMarketplaceGetPageRewardReq) GetNeedSearch() interface{} {
	return *m
}

type MarketplaceValidationUpdateReq struct {
	Validations []struct {
		Id    int  `json:"id"`
		APass bool `json:"a_pass"`
	} `json:"validations"`
}

func (s *MarketplaceValidationUpdateReq) GetId() interface{} {
	return s.Validations[0].Id
}

type DataHubMarketplaceGetCampaignDisputeReq struct {
	dto.Pagination `search:"-"`
	TaskID         int `form:"task_id" search:"type:exact;column:task;table:ai_task_upload_records"`
	UID            int `form:"uid" search:"type:exact;column:user;table:ai_task_upload_records"`
}

func (m *DataHubMarketplaceGetCampaignDisputeReq) GetNeedSearch() interface{} {
	return *m
}

type MarketplaceDisputeUpdateReq struct {
	TaskID      int `json:"task_id"`
	Validations []struct {
		Id    int  `json:"id"`
		APass bool `json:"a_pass"`
	} `json:"validations"`
}

func (s *MarketplaceDisputeUpdateReq) GetId() interface{} {
	return s.Validations[0].Id
}

type MarketplaceValidationMaliciousUpdateReq struct {
	Id   int  `json:"id"`
	Flag bool `json:"m_flag"`
}

func (s *MarketplaceValidationMaliciousUpdateReq) GetId() interface{} {
	return s.Id
}

type GetCampaignValidationSummaryReq struct {
	TaskID    int    `form:"task_id" search:"type:exact;column:task;table:ai_task_upload_records"`
	StartTime string `form:"start_time" search:"type:gte;column:created_at;table:ai_task_upload_records"`
	EndTime   string `form:"end_time" search:"type:lte;column:created_at;table:ai_task_upload_records"`
	UID       int    `form:"uid" search:"type:exact;column:user;table:ai_task_upload_records"`
	Status    string `form:"status" search:"type:exact;column:status;table:ai_task_upload_records"`
}

type AddCampaignReq struct {
	Name               string `json:"name"`
	Cover              string `json:"cover"`
	Level              int    `json:"level"` // 0 Easy 1 Medium 2 Difficult 3 Extremely Hard
	Start              int    `json:"start"`
	End                int    `json:"end"`
	Type               int    `json:"type"` // 0 Image 1 Text 2 Video 3 Audio
	Tags               string `json:"tags"` // ["Data collection"]
	Description        string `json:"description"`
	SimpleDescription  string `json:"simple_description"`
	Requirements       string `json:"requirements"`
	USDTReward         string `json:"usdt_reward"`
	PointReward        string `json:"point_reward"`
	MaxSize            int    `json:"max_size"`   // max size of file, byte,0 unlimited
	MaxNumber          int    `json:"max_number"` // max number of files
	MinNumber          int    `json:"min_number"` // min number of files
	Conditions         string `json:"conditions"` // ["email", "solana"]
	PointStake         string `json:"point_stake"`
	VerifyRequirements string `json:"verify_requirements"`
}

type UpdateCampaignReq struct {
	TaskID             int    `json:"task_id"`
	Name               string `json:"name"`
	Cover              string `json:"cover"`
	Level              int    `json:"level"` // 0 Easy 1 Medium 2 Difficult 3 Extremely Hard
	Start              int    `json:"start"`
	End                int    `json:"end"`
	Type               int    `json:"type"` // 0 Image 1 Text 2 Video 3 Audio
	Tags               string `json:"tags"` // ["Data collection"]
	Description        string `json:"description"`
	SimpleDescription  string `json:"simple_description"`
	Requirements       string `json:"requirements"`
	USDTReward         string `json:"usdt_reward"`
	PointReward        string `json:"point_reward"`
	MaxSize            int    `json:"max_size"`   // max size of file, byte,0 unlimited
	MaxNumber          int    `json:"max_number"` // max number of files
	MinNumber          int    `json:"min_number"` // min number of files
	Conditions         string `json:"conditions"` // ["email", "solana"]
	PointStake         string `json:"point_stake"`
	VerifyRequirements string `json:"verify_requirements"`
}

type DeleteCampaignReq struct {
	TaskID int `json:"task_id"`
}

type CampaignUploadReq struct {
	Files []multipart.FileHeader `form:"files" swaggerignore:"true"`
}

type CampaignUploadResponse struct {
	Links []string `json:"links"`
}

type CampaignDetailReq struct {
	Id int `uri:"id"`
}

type CampaignDetailResponse struct {
	TaskID             int    `json:"task_id"`
	Name               string `json:"name"`
	Cover              string `json:"cover"`
	Level              int    `json:"level"` // 0 Easy 1 Medium 2 Difficult 3 Extremely Hard
	Start              int    `json:"start"`
	End                int    `json:"end"`
	Type               int    `json:"type"` // 0 Image 1 Text 2 Video 3 Audio
	Tags               string `json:"tags"` // ["Data collection"]
	Description        string `json:"description"`
	SimpleDescription  string `json:"simple_description"`
	Requirements       string `json:"requirements"`
	USDTReward         string `json:"usdt_reward"`
	PointReward        string `json:"point_reward"`
	MaxSize            int    `json:"max_size"`   // max size of file, byte,0 unlimited
	MaxNumber          int    `json:"max_number"` // max number of files
	MinNumber          int    `json:"min_number"` // min number of files
	Conditions         string `json:"conditions"` // ["email", "solana"]
	PointStake         string `json:"point_stake"`
	VerifyRequirements string `json:"verify_requirements"`
}
