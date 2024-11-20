package dto

import "go-admin/common/dto"

type DataHubUserGetPageReq struct {
	dto.Pagination `search:"-"`
	Email          string `form:"email" search:"type:contains;column:email;table:users"`
	TelegramName   string `form:"telegram_name" search:"type:contains;column:telegram_name;table:users"`
	TwitterName    string `form:"twitter_name" search:"type:contains;column:twitter_name;table:users"`
	SolanaAccount  string `form:"solana_account" search:"type:contains;column:wallet;table:users"`
	DataHubUserOrder
}

type DataHubUserOrder struct {
	RankOrder string `search:"type:order;column:rank;table:users" form:"rankOrder"`
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
