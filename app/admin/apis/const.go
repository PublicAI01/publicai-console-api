package apis

const (
	SENIOR_POINT = 40000
	MASTER_POINT = 120000
)

type UserLevelType uint

const (
	Beginner UserLevelType = iota + 1
	Senior
	Master
)
