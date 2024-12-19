package models

import "time"

type DataHubUser struct {
	ID uint `gorm:"primaryKey;type:autoIncrement;autoIncrementIncrement:1" json:"id"`
	//Identity         string    `gorm:"not null;type:varchar(20)" json:"identity"`
	Email *string `gorm:"type:varchar(40);default:null;index" json:"email"`
	Name  string  `gorm:"type:varchar(20);default:null" json:"name"`
	//Avatar           *string   `gorm:"type:varchar(300);default:null" json:"avatar"`
	Wallet *string `gorm:"type:varchar(100);default:null" json:"wallet"`
	//Nonce            *string   `gorm:"type:varchar(10);default:null" json:"nonce"`
	//InviteCode       string    `gorm:"type:varchar(10);default:null" json:"invite_code"`
	Point uint `json:"point"`
	//TwitterID   *string `gorm:"type:varchar(40);default:null" json:"twitter_id"`
	TwitterName *string `gorm:"type:varchar(20);default:null" json:"twitter_name"`
	//FollowersCount   uint      `gorm:"default:0" json:"followers_count"`
	Level uint `gorm:"default:1" json:"level"` //1 Beginner 2 Senior 3 Master
	//SucCount         uint      `gorm:"default:0" json:"suc_count"`
	//DoneCount uint `gorm:"default:0" json:"done_count"`
	//LoginIP          *string   `json:"login_ip"`
	Location *string `json:"location"`
	//SBTLevel         int       `gorm:"default:-1" json:"sbt_level"` //-1 none 1 Beginner 2 Senior 3 Master
	//SBTTxHash *string `json:"sbt_tx_hash"`
	//BuilderPoint     int       `json:"builder_point"`
	TelegramID        *string   `gorm:"type:varchar(40);default:null" json:"-"`
	TelegramName      *string   `gorm:"type:varchar(40);default:null" json:"telegram_name"`
	TelegramFullName  *string   `gorm:"type:varchar(40);default:null" json:"-"`
	NearAccount       *string   `json:"near_account"`
	EvmAccount        *string   `json:"evm_account"`
	Rank              int       `json:"rank"`
	CompletedItems    int       `json:"completed_items"`
	UploadTimes       int       `json:"upload_times"`
	ContributionValue int       `json:"contribution_value"`
	CreatedAt         time.Time `gorm:"type:timestamptz;autoCreateTime:milli;default:null" json:"created_at"`
}

func (*DataHubUser) TableName() string {
	return "users"
}

type Train struct {
	ID           uint      `gorm:"primaryKey;type:autoIncrement;autoIncrementIncrement:1" json:"id"`
	User         uint      `gorm:"index" json:"user"`
	Dataset      uint      `gorm:"index" json:"dataset"`
	Selected     int       `json:"selected"`
	TrainJob     uint      `gorm:"index" json:"train_job"`
	DatasetType  uint      `json:"dataset_type"`
	Valid        bool      `gorm:"default:false" json:"valid"`
	AudioContent *string   `gorm:"default:null" json:"audio_content"`
	AudioSimilar int       `gorm:"default:0" json:"audio_similar"`
	CreatedAt    time.Time `gorm:"type:timestamptz;autoCreateTime:milli;default:null" json:"created_at"`
}

func (*Train) TableName() string {
	return "trains"
}

type TMAUser struct {
	ID               uint      `gorm:"primaryKey;type:autoIncrement;autoIncrementIncrement:1" json:"id"`
	AppID            int       `gorm:"default:0" json:"app_id"`
	TelegramID       string    `gorm:"index" json:"telegram_id"`
	TelegramFullName string    `gorm:"type:varchar(40);default:null" json:"telegram_full_name"`
	TMAPoint         int       `json:"tma_point"`
	InviteCode       string    `gorm:"type:varchar(10);default:null;index" json:"invite_code"`
	CreatedAt        time.Time `gorm:"type:timestamptz;autoCreateTime:milli;default:null" json:"created_at"`
}

func (*TMAUser) TableName() string {
	return "tma_users"
}

type UserNearAddress struct {
	ID          uint      `gorm:"primaryKey;type:autoIncrement;autoIncrementIncrement:1" json:"id"`
	User        uint      `json:"user"`
	NearAddress *string   `gorm:"default:null" json:"near_address"`
	CreatedAt   time.Time `gorm:"type:timestamptz;autoCreateTime:milli;default:null" json:"created_at"`
}

func (*UserNearAddress) TableName() string {
	return "user_near_addresses"
}

type UserEthereumAddress struct {
	ID              uint      `gorm:"primaryKey;type:autoIncrement;autoIncrementIncrement:1" json:"id"`
	User            uint      `json:"user"`
	EthereumAddress *string   `gorm:"type:varchar(50);default:null" json:"ethereum_address"`
	CreatedAt       time.Time `gorm:"type:timestamptz;autoCreateTime:milli;default:null" json:"created_at"`
}

func (*UserEthereumAddress) TableName() string {
	return "user_ethereum_addresses"
}

type TrainReward struct {
	ID        uint      `gorm:"primaryKey;type:autoIncrement;autoIncrementIncrement:1" json:"id"`
	User      uint      `gorm:"index" json:"user"`
	TrainJob  int       `json:"train_job"` //0 train -1 bind wallet -2 invite -3 bind email -4 collection -5 bug feedback -6 invite rebate -7 DBReport -8 DBVerifiedX -9 DBZealyReward -10 DBGPTHunterReward
	Dataset   uint      `json:"dataset"`
	Point     uint      `json:"point"`
	Type      uint      `gorm:"default:0" json:"type"` // 0 train 1 bind wallet 2 invite 3 bind email 4 collection 5 bug feedback 6 invite rebate 7 report 8 verifiedX 9 ZealyReward 10 GPTHunterReward
	CreatedAt time.Time `gorm:"type:timestamptz;autoCreateTime:milli;default:null" json:"created_at"`
}

func (*TrainReward) TableName() string {
	return "train_rewards"
}

type RewardItem struct {
	ID        int       `json:"id"`
	User      uint      `json:"user"`
	Point     int       `json:"point"`
	CreatedAt time.Time `json:"created_at"`
}

type AllRewardItem struct {
	User  uint `json:"user"`
	Point uint `json:"point"`
}

type DataHubAmbassador struct {
	ID                    uint    `gorm:"primaryKey;type:autoIncrement;autoIncrementIncrement:1" json:"id"`
	Email                 *string `gorm:"type:varchar(40);default:null;index" json:"email"`
	Name                  string  `gorm:"type:varchar(20);default:null" json:"name"`
	TwitterName           *string `gorm:"type:varchar(20);default:null" json:"twitter_name"`
	Location              *string `json:"location"`
	TelegramID            *string `gorm:"type:varchar(40);default:null" json:"-"`
	TelegramName          *string `gorm:"type:varchar(40);default:null" json:"telegram_name"`
	TelegramFullName      *string `gorm:"type:varchar(40);default:null" json:"-"`
	ConsensusContribution int     `json:"consensus_contribution"`
}

func (*DataHubAmbassador) TableName() string {
	return "users"
}
