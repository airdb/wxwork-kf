package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NICEXAI/WeChatCustomerServiceSDK/syncmsg"
	"github.com/airdb/wxwork-kf/internal/app"
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
			OriginData: generateTextData("帮助"),
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
