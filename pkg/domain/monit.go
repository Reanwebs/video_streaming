package domain

type VideoRewardRequest struct {
	UserID    string
	VideoID   uint32
	Reason    string
	Views     uint32
	PaidCoins uint32
}
