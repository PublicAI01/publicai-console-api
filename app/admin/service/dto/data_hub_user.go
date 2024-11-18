package dto

import "go-admin/common/dto"

//type DataHubUser struct {
//	UserID            string `json:"user_id"`
//	UserName          string `json:"user_name"`
//	Location          string `json:"location"`
//	Email             string `json:"email"`
//	TelegramID        string `json:"telegram_id"`
//	TelegramName      string `json:"telegram_id"`
//	TwitterID         string `json:"twitter_id"`
//	TwitterName       string `json:"twitter_name"`
//	Level             int    `json:"level"`
//	Rank              int    `json:"rank"`
//	SolanaAccount     string `json:"solana_account"`
//	NearAccount       string `json:"near_account"`
//	EvmAccount        string `json:"evm_account"`
//	CompletedItems    int    `json:"completed_items"`
//	Points            int    `json:"points"`
//	UploadTimes       int    `json:"upload_times"`
//	ContributionValue int    `json:"contribution_value"`
//}

type DataHubUserGetPageReq struct {
	dto.Pagination `search:"-"`
	Email          string `form:"email" search:"type:contains;column:email;table:users"`
	TelegramName   string `form:"telegram_name" search:"type:contains;column:telegram_name;table:users"`
	TwitterName    string `form:"twitter_name" search:"type:contains;column:twitter_name;table:users"`
	SolanaAccount  string `form:"solana_account" search:"type:contains;column:wallet;table:users"`
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
