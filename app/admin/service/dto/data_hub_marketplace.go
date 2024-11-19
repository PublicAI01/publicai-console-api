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
	TaskID         int `form:"task_id" search:"type:exact;column:task;table:ai_task_uploaded_files"`
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
	Id    int  `uri:"id"`
	APass bool `json:"a_pass"`
}

func (s *MarketplaceValidationUpdateReq) GetId() interface{} {
	return s.Id
}
