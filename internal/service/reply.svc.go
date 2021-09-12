package service

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	sdk "github.com/NICEXAI/WeChatCustomerServiceSDK"
	"github.com/NICEXAI/WeChatCustomerServiceSDK/sendmsg"
	"github.com/NICEXAI/WeChatCustomerServiceSDK/syncmsg"
	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/airdb/wxwork-kf/pkg/po"
	"github.com/airdb/wxwork-kf/pkg/util"
)

type Reply struct{}

func NewReply() *Reply {
	return &Reply{}
}

// ProcMsg 处理单条消息, 并按消息来源颁发给不同的处理过程
func (s Reply) ProcMsg(ctx context.Context, msg syncmsg.Message) {
	switch msg.Origin {
	case 3: //客户回复的消息
		s.userMsg(ctx, msg)
	case 4: //系统推送的消息
		s.systemMsg(ctx, msg)
	default:
		log.Fatalf("unknown msg origin: %d", msg.Origin)
	}
}

// 处理客户回复的消息
func (s Reply) userMsg(ctx context.Context, msg syncmsg.Message) {
	// 按整块消息进行匹配
	rTpl, ok := TplReplys.Match(msg)
	if !ok {
		rTpl, _ = TplReplys.Default()
	}

	// 最终的冗余，这块代码应该不被执行
	if rTpl == nil {
		rTpl = &ReplyTpl{
			ReplyType: ReplyTypeText,
			Message:   DefaultMsg,
		}
	}

	var (
		ret          = rTpl.Gen(msg.ExternalUserID, msg.OpenKFID, util.RandString(32))
		hasMsgSendOk bool // 消息执行是否成功
	)

	switch rTpl.ReplyType {
	case ReplyTypeText, ReplyTypeImage, ReplyTypeMenu:
		hasMsgSendOk = s.sendMsg(msg, ret)
	case ReplyTypeActionTrans: // 分配客服会话
		hasMsgSendOk = s.transMsg(msg, ret)
	}

	if hasMsgSendOk {
		s.saveMsg()
	}
}

// 处理系统消息
func (s Reply) systemMsg(ctx context.Context, msg syncmsg.Message) {
	var (
		ret          interface{}
		sentCackeKey string // 消息缓存key
		sentCackeTTL time.Duration
	)

	switch msg.EventType {
	case "enter_session": // 用户进入会话事件
		tMsg, _ := msg.GetEnterSessionEvent()
		// 缓存上次该客户的欢迎消息发送，避免重复发送。
		sentCackeKey = strings.Join([]string{
			app.SentMsgPrefix, msg.EventType, tMsg.OpenKFID, tMsg.ExternalUserID,
		}, ":")
		sentCackeTTL = 6 * time.Hour
		// 检查最近6小时是否发送过
		lastSend, _ := app.Redis.Get(ctx, sentCackeKey).Result()
		if len(lastSend) > 0 {
			return
		}

		rMsg := &sendmsg.Text{
			Message: sendmsg.Message{
				ToUser:   tMsg.ExternalUserID,
				OpenKFID: tMsg.OpenKFID,
				MsgID:    util.RandString(32),
			},
			MsgType: "text",
		}
		rMsg.Text.Content = WelcomeMsg // 欢迎语
		ret = rMsg
	case "msg_send_fail": // 消息发送失败事件
		fallthrough
	case "servicer_status_change": // 客服人员接待状态变更事件
		fallthrough
	case "session_status_change": // 会话状态变更事件
		fallthrough
	default:
		log.Fatalf("unknown system event type: %s", msg.EventType)
	}

	if s.sendMsg(msg, ret) && len(sentCackeKey) > 0 {
		app.Redis.Set(ctx, sentCackeKey, time.Now().String(), sentCackeTTL).Result()
	}
}

// 发送消息
func (s Reply) sendMsg(msg syncmsg.Message, ret interface{}) bool {
	if info, err := app.WxClient.SendMsg(ret); err == nil {
		params, _ := json.Marshal(ret)
		log.Println("result:", msg.EventType, info.MsgID, ", msg:", string(params))

		return true
	} else {
		log.Println("result:", msg.EventType, ", err:", err.Error())

		return false
	}
}

// 分配客服会话
func (s Reply) transMsg(msg syncmsg.Message, ret interface{}) bool {
	transMsg, ok := ret.(sdk.ServiceStateTransOptions)
	if !ok {
		return false
	}
	transInfo, err := app.WxClient.ServiceStateTrans(transMsg)
	if err != nil {
		log.Fatalf("trans msg err(%d): %s", transInfo.ErrCode, transInfo.ErrMsg)
		return false
	}

	return true
}

// 执行消息持久化
func (s Reply) saveMsg() {
	var (
		userInput         string
		wxResponseContent string
		wxMessageId       string
		wxOpenKFID        string
		toUser            string
	)

	logData := new(po.WxKfLog)
	logData.MsgID = wxMessageId
	logData.OpenKFID = wxOpenKFID
	logData.ToUserID = toUser
	logData.Input = userInput
	logData.Response = wxResponseContent
	po.WxKfLogSave(logData)
}
