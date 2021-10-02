package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NICEXAI/WeChatCustomerServiceSDK/syncmsg"
	"github.com/airdb/sailor/dbutil"
	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/airdb/wxwork-kf/internal/store/mysql"
	"github.com/airdb/wxwork-kf/internal/types"
	"github.com/airdb/wxwork-kf/pkg/schema"
)

func TestReply_ProcMsg(t *testing.T) {
	app.InitSdk()

	type args struct {
		ctx context.Context
		msg syncmsg.Message
	}

	tests := []struct {
		name string
		s    *Reply
		args args
	}{
		// {"text msg", NewReply(), args{context.Background(), syncmsg.Message{
		// 	Origin:     3,
		// 	MsgType:    "text",
		// 	OriginData: generateTextData("[寻人]"),
		// }}},
		{"menu msg", NewReply(nil), args{context.Background(), syncmsg.Message{
			Origin:     3,
			MsgType:    "text",
			OriginData: generateTextData("志愿者"),
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Reply{}
			s.ProcMsg(tt.args.ctx, tt.args.msg)
		})
	}
}

func generateTextData(s string) []byte {
	baseMsg := syncmsg.BaseMessage{
		Origin: 3,
	}
	msg := syncmsg.Text{
		BaseMessage: baseMsg,
	}
	msg.MsgType = "text"
	msg.Text.Content = s

	b, _ := json.Marshal(msg)

	return b
}

func TestReply_saveMsg(t *testing.T) {
	dbutil.InitDefaultDB()
	db := dbutil.WriteDefaultDB()
	db.Migrator().AutoMigrate(&schema.Talk{}, &schema.Message{})

	store, _ := mysql.GetFactoryOr(db)
	reply := NewReply(store)
	type args struct {
		ctx  context.Context
		data interface{}
	}
	tests := []struct {
		name string
		s    Reply
		args args
	}{
		{"", *reply, args{context.Background(), &types.ReplyMessage{
			OpenKFID:    "test_open_kfid",
			ToUser:      "test_to_user",
			ReplyType:   types.WxMsgTypeText,
			ContentText: "test_content_text",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.saveMsg(tt.args.ctx, tt.args.data)
		})
	}
}
