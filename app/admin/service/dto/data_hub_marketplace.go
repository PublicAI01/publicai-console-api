package dto

import (
	"go-admin/common/dto"
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
	TaskID         int `form:"task_id" search:"type:exact;column:task;table:ai_task_upload_records"`
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
