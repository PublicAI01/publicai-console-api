package models

import "time"

type FileItem struct {
	ID    int    `json:"id"`
	Link  string `json:"link"`
	VPass bool   `json:"v_pass"`
	APass bool   `json:"a_pass"`
}
type AITaskShowRecordItem struct {
	NO              int        `json:"no"`
	ID              uint       `json:"id"`
	DataNumber      int        `json:"data_number"`
	User            uint       `json:"user"`
	Status          int        `json:"status"`
	Editor          int        `json:"editor"` // -1 nil 0 管理员 1 AI
	Valid           int        `json:"valid"`
	UploadTime      time.Time  `json:"upload_time"`
	Items           []FileItem `json:"items" gorm:"-"`
	MaliciousReason int        `json:"malicious_reason"` // 0 nil 1 miss target 2 AI generated
	//Items []AITaskUploadedFile `json:"items" gorm:"-"`
}

type FileDisputeItem struct {
	ID   int    `json:"id"`
	Link string `json:"link"`
	VAye int    `json:"v_aye"`
	VNay int    `json:"v_nay"`
}
type AITaskShowDisputeRecordItem struct {
	NO         int               `json:"no"`
	ID         uint              `json:"id"`
	DataNumber int               `json:"data_number"`
	User       uint              `json:"user"`
	Status     int               `json:"status"`
	Valid      int               `json:"valid"`
	UploadTime time.Time         `json:"upload_time"`
	Items      []FileDisputeItem `json:"items" gorm:"-"`
}

type AITaskUploadRecord struct {
	ID         uint                 `gorm:"primaryKey;type:autoIncrement;autoIncrementIncrement:1" json:"id"`
	User       uint                 `gorm:"index" json:"user"`
	Task       uint                 `json:"task"`
	TaskRecord uint                 `gorm:"index" json:"task_record"`
	Type       int                  `json:"type"` // 0 Image 1 Text 2 Video 3 Audio
	Total      int                  `json:"total"`
	Success    int                  `json:"success"`
	Files      string               `json:"files"`
	Issued     int                  `gorm:"default:0" json:"issued"`
	Completed  int                  `gorm:"default:0" json:"completed"`
	CanIssue   bool                 `gorm:"default:false" json:"can_issue"`
	Status     int                  `gorm:"default:0" json:"status"` // -1 Subpar(单次上传全部不合格) 0 Pending Verification 1 Partially Accepted 2 All Accepted
	CreatedAt  time.Time            `gorm:"type:timestamptz;autoCreateTime:milli;default:null" json:"created_at"`
	Items      []AITaskUploadedFile `json:"items" gorm:"-"`
}

func (*AITaskUploadRecord) TableName() string {
	return "ai_task_upload_records"
}

type AITaskUploadedFile struct {
	ID           uint      `gorm:"primaryKey;type:autoIncrement;autoIncrementIncrement:1" json:"id"`
	User         uint      `gorm:"index" json:"user"`
	Task         uint      `gorm:"index:idx_task_hash" json:"task"`
	TaskRecord   uint      `gorm:"index" json:"task_record"`
	UploadRecord uint      `gorm:"index" json:"upload_record"`
	Type         int       `json:"type"` // 0 Image 1 Text 2 Video 3 Audio
	Link         string    `json:"link"`
	Hash         string    `gorm:"index:idx_task_hash" json:"hash"`
	FileName     string    `gorm:"index" json:"file_name"`
	VPass        bool      `gorm:"default:false" json:"v_pass"` // validator verify passed
	APass        bool      `gorm:"default:true" json:"a_pass"`  // admin verify passed
	Status       int       `gorm:"default:0" json:"status"`     //最终状态 -1 Data not qualified for task 0 accepted 1 finished
	VAye         int       `gorm:"default:0" json:"v_aye"`      //aye of validator
	VNay         int       `gorm:"default:0" json:"v_nay"`      //nay of validator
	CreatedAt    time.Time `gorm:"type:timestamptz;autoCreateTime:milli;default:null" json:"created_at"`
}

func (*AITaskUploadedFile) TableName() string {
	return "ai_task_uploaded_files"
}

type AITask struct {
	ID                 uint      `gorm:"primaryKey;type:autoIncrement;autoIncrementIncrement:1" json:"id"`
	Name               string    `gorm:"type:varchar(100);default:null;index" json:"name"`
	Level              int       `json:"level"` // 0 Easy 1 Medium 2 Difficult 3 Extremely Hard
	Start              time.Time `gorm:"type:timestamptz" json:"start"`
	End                time.Time `gorm:"type:timestamptz" json:"end"`
	Type               int       `json:"type"` // 0 Image 1 Text 2 Video 3 Audio
	Tags               string    `gorm:"type:varchar(100);default:null;index" json:"tags"`
	Description        string    `gorm:"type:varchar(4000);default:null;index" json:"description"`
	SimpleDescription  string    `gorm:"type:varchar(1000);default:null;" json:"simple_description"`
	Requirements       string    `gorm:"type:varchar(4000);default:null" json:"requirements"`
	Examples           string    `json:"examples"`
	IneligibleExamples string    `json:"ineligible_examples"`
	RewardPool         string    `gorm:"type:varchar(100);default:null" json:"reward_pool"`
	FormLink           string    `gorm:"type:varchar(300);default:null" json:"form_link"`
	Collected          int       `json:"collected"`
	Required           int       `json:"required"`
	USDTReward         int       `json:"usdt_reward"`
	PointReward        int       `json:"point_reward"`
	MaxSize            int       `json:"max_size"`   // max size of file, byte,0 unlimited
	MaxNumber          int       `json:"max_number"` // max number of files
	MinNumber          int       `json:"mn_number"`  // min number of files
	Conditions         string    `json:"conditions"` // ["email", "solana"]
	Consensus          int       `gorm:"default:201" json:"consensus"`
	MinValid           int       `gorm:"default:5" json:"min_valid"`
	Available          bool      `gorm:"default:true" json:"available"`
	USDTStake          int       `gorm:"default:0" json:"usdt_stake"`
	PointStake         int       `gorm:"default:0" json:"point_stake"`
	TelegramGroup      *string   `json:"telegram_group"`
	CheckForReward     bool      `gorm:"default:false" json:"check_for_reward"`
	CreatedAt          time.Time `gorm:"type:timestamptz;autoCreateTime:milli;default:null" json:"created_at"`
}

func (*AITask) TableName() string {
	return "ai_tasks"
}

type AITaskRecord struct {
	ID          uint      `gorm:"primaryKey;type:autoIncrement;autoIncrementIncrement:1" json:"id"`
	Task        uint      `json:"task"`
	User        uint      `json:"user"`
	Status      int       `json:"status"` //-2 Data not qualified for task -1 Didn't complete the task 0 accepted 1 finished
	Number      int       `json:"number"`
	USDTReward  int       `json:"usdt_reward"`
	PointReward int       `json:"point_reward"`
	USDTLink    string    `json:"usdt_link"`
	PointLink   string    `json:"point_link"`
	UpdateAt    time.Time `gorm:"type:timestamptz;autoCreateTime:milli;default:null" json:"update_at"`
	CreatedAt   time.Time `gorm:"type:timestamptz;autoCreateTime:milli;default:null" json:"created_at"`
}

func (*AITaskRecord) TableName() string {
	return "ai_task_records"
}

type MarketPlaceRewardItem struct {
	No            uint      `json:"no"`
	User          uint      `json:"user"`
	Email         *string   `json:"email"`
	SolanaAccount *string   `json:"solana_account"`
	USDTReward    string    `json:"usdt_reward"`
	PointReward   string    `json:"point_reward"`
	Collected     int       `json:"collected"`
	Accepted      int       `json:"accepted"`
	Indicators    int       `json:"indicators"`
	USDTLink      *string   `json:"usdt_link"`
	CreatedAt     time.Time `json:"created_at"`
}

type AITaskVariants struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Start            time.Time `json:"start"`
	End              time.Time `json:"end"`
	UploadTimes      int       `json:"upload_times"`
	AmountOfAll      int       `json:"amount_of_all"`
	TotalUploadUsers int       `json:"total_upload_users"`
	AcceptedTimes    int       `json:"accepted_times"`
	AcceptedTotal    int       `json:"accepted_total"`
	AcceptedUsers    int       `json:"accepted_users"`
}

type ValidationSummary struct {
	AllAccepted int `json:"all_accepted"`
	//AllAcceptedTotal     int `json:"all_accepted_total"`
	PartialAccepted int `json:"partial_accepted"`
	//PartialAcceptedTotal int `json:"partial_accepted_total"`
	Subpar int `json:"subpar"`
	//SubparTotal          int `json:"subpar_total"`
	Malicious int `json:"malicious"`
	//MaliciousTotal       int `json:"malicious_total"`
	Pending int `json:"pending"`
	//PendingTotal         int `json:"pending_total"`
	Total             int `json:"total"`
	AmbassadorPending int `json:"ambassador_pending"`
}

func (*ValidationSummary) TableName() string {
	return "ai_task_upload_records"
}
