package domain

type VideoRewardRequest struct {
	UserID    string
	VideoID   string
	Reason    string
	Views     uint32
	PaidCoins uint32
}

type ExclusiveContentRequest struct {
	UserID    string
	VideoID   string
	Reason    string
	Owner     string
	PaidCoins uint32
}
