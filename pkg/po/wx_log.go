package po

import (
	"fmt"
	"github.com/airdb/sailor/dbutil"
	"gorm.io/gorm"
	"time"
)

type WxKfLog struct {
	gorm.Model
	Input    string `gorm:"type:text" json:"input"`
	Response string `gorm:"type:text" json:"response"`
	OpenKFID string `gorm:"type:varchar(100);column:open_kfid" json:"open_kfid"`
	ToUserID string `gorm:"type:varchar(100)" json:"to_user_id"`
	MsgID    string `gorm:"type:varchar(100)"json:"msg_id"`
}

func (WxKfLog) TableName() string {
	return "tab_wx_kf_log"
}

func (u *WxKfLog) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
func (u *WxKfLog) BeforeSave(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return nil
}

func WxKfLogSave(data *WxKfLog) error {
	fmt.Println(data)
	return dbutil.WriteDefaultDB().Debug().Create(data).Error
}
