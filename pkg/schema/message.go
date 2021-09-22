package schema

import "time"

type Message struct {
	ModelMeta
	TalkID         ModelPk   `gorm:"type:varchar(100);column:talk_id"        json:"talkId"`
	MsgFrom        string    `gorm:"type:varchar(16);column:msg_from"        json:"msgFrom"`
	Origin         uint32    `gorm:"type:tinyint;column:origin"              json:"origin"`
	Msgid          string    `gorm:"type:varchar(100);column:msgid"          json:"msgid"`
	Msgtype        string    `gorm:"type:varchar(16);column:msgtype"         json:"msgtype"`
	SendTime       time.Time `gorm:"type:timestamp;column:send_time"         json:"sendTime"`
	ServicerUserid string    `gorm:"type:varchar(100);column:service_userid" json:"serviceUserid"`
	Content        string    `gorm:"type:text;column:content"                json:"content"`
	Raw            string    `gorm:"type:text;column:raw"                    json:"raw"`
}
