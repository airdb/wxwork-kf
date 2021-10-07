package schema

type Talk struct {
	ModelMeta
	OpenKFID string `gorm:"type:varchar(100);column:open_kfid" json:"openKfid"`
	ToUserID string `gorm:"type:varchar(100);column:to_user_id" json:"toUserId"`
}
