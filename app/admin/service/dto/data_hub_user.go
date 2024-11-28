package dto

import "go-admin/common/dto"

type DataHubUserGetPageReq struct {
	dto.Pagination `search:"-"`
	ID             int    `form:"id" search:"type:exact;column:id;table:users"`
	UserName       string `form:"name" search:"type:exact;column:name;table:users"`
	Email          string `form:"email" search:"type:exact;column:email;table:users"`
	TelegramName   string `form:"telegram_name" search:"type:exact;column:telegram_name;table:users"`
	TwitterName    string `form:"twitter_name" search:"type:exact;column:twitter_name;table:users"`
	SolanaAccount  string `form:"solana_account" search:"type:exact;column:wallet;table:users"`
	DataHubUserOrder
}

type DataHubUserOrder struct {
	RankOrder      string `search:"type:order;column:rank;table:users" form:"rankOrder"`
	PointOrder     string `search:"type:order;column:point;table:users" form:"pointOrder"`
	CompletedOrder string `search:"type:order;column:point;table:users" form:"completed_itemsOrder"`
}

func (m *DataHubUserGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type DataHubUserPointGetPageReq struct {
	dto.Pagination `search:"-"`
	UID            int    `form:"uid" search:"type:exact;column:user;table:trains"`
	StartTime      string `form:"start_time" search:"type:gte;column:created_at;table:trains"`
	EndTime        string `form:"end_time" search:"type:lte;column:created_at;table:trains"`
}

func (m *DataHubUserPointGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type DataHubUserGetAllRewardReq struct {
	dto.Pagination `search:"-"`
	StartTime      string `form:"start_time" search:"type:gte;column:created_at;table:train_rewards"`
	EndTime        string `form:"end_time" search:"type:lte;column:created_at;table:train_rewards"`
	DataHubAllRewardOrder
}

type DataHubAllRewardOrder struct {
	PointOrder string `search:"type:order;column:point;table:users" form:"pointOrder"`
}

func (m *DataHubUserGetAllRewardReq) GetNeedSearch() interface{} {
	return *m
}

type DataHubAmbassadorGetPageReq struct {
	dto.Pagination `search:"-"`
	ID             int `form:"id" search:"type:exact;column:id;table:users"`
	DataHubAmbassadorOrder
}

type DataHubAmbassadorOrder struct {
	ContributionOrder string `search:"type:order;column:consensus_contribution;table:users" form:"consensus_contributionOrder"`
}

func (m *DataHubAmbassadorGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type MarketplaceAmbassadorUpdateReq struct {
	Ambassadors []struct {
		Id   int  `json:"id"`
		Flag bool `json:"flag"`
	} `json:"ambassadors"`
}

func (s *MarketplaceAmbassadorUpdateReq) GetId() interface{} {
	return s.Ambassadors[0].Id
}

type DataHubExportAmbassadorReq struct {
	dto.Pagination `search:"-"`
	DataHubExportAmbassadorOrder
}

type DataHubExportAmbassadorOrder struct {
	ContributionOrder string `search:"type:order;column:point;table:users" form:"contributionOrder"`
}

func (m *DataHubExportAmbassadorReq) GetNeedSearch() interface{} {
	return *m
}
