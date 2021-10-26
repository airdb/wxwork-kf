package service

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/airdb/wxwork-kf/internal/store"
	"github.com/airdb/wxwork-kf/internal/types"
	"github.com/airdb/wxwork-kf/pkg/schema"
	"github.com/airdb/wxwork-kf/pkg/util"
	"github.com/silenceper/wechat/v2/work/kf"
	"github.com/silenceper/wechat/v2/work/kf/sendmsg"
	"github.com/silenceper/wechat/v2/work/kf/syncmsg"
)

type Reply struct {
	store store.Factory // TODO
}

func NewReply(store store.Factory) *Reply {
	return &Reply{store}
}

// ProcMsg 处理单条消息, 并按消息来源颁发给不同的处理过程
func (s Reply) ProcMsg(ctx context.Context, msg syncmsg.Message) { // 记录用户发送的消息
	s.saveMsg(ctx, &msg)
	switch msg.Origin {
	case 3: // 客户回复的消息
		s.userMsg(ctx, msg)
	case 4: // 系统推送的消息
		s.systemMsg(ctx, msg)
	case 5: // 接待人员在企业微信客户端发送的消息
		s.receptionistMsg(ctx, msg)
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
		msgResp      = rTpl.Gen(msg.ExternalUserID, msg.OpenKFID)
		hasMsgSendOk bool // 消息执行是否成功
	)

	ret, err := msgResp.Assume()
	if err != nil {
		log.Fatalf("can not assume msg: %s", err.Error())
		return
	}

	switch rTpl.ReplyType {
	case ReplyTypeText, ReplyTypeImage, ReplyTypeMenu:
		hasMsgSendOk = s.sendMsg(msg, ret)
	case ReplyTypeActionTrans: // 分配客服会话
		hasMsgSendOk = s.transMsg(msg, ret)
	}

	if hasMsgSendOk {
		s.saveMsg(ctx, msgResp)
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

// 处理客服消息, 只需用入库
func (s Reply) receptionistMsg(ctx context.Context, msg syncmsg.Message) {
	s.saveMsg(ctx, msg)
}

// 发送消息
func (s Reply) sendMsg(msg syncmsg.Message, ret interface{}) bool {
	if info, err := app.WxWorkKF.SendMsg(ret); err == nil {
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
	transMsg, ok := ret.(kf.ServiceStateTransOptions)
	if !ok {
		return false
	}
	transInfo, err := app.WxWorkKF.ServiceStateTrans(transMsg)
	if err != nil {
		log.Fatalf("trans msg err(%d): %s", transInfo.ErrCode, transInfo.ErrMsg)
		return false
	}

	return true
}

// 执行消息持久化
// TODO: 根据消息内容执行不同的持久化方式
func (s Reply) saveMsg(ctx context.Context, data interface{}) {
	var (
		talk *schema.Talk
		msg  *schema.Message
	)

	switch m := data.(type) {
	case *types.ReplyMessage: // 返回的消息
		talk = &schema.Talk{
			OpenKFID: m.OpenKFID,
			ToUserID: m.ToUser,
		}
		msg = &schema.Message{
			MsgFrom:  "bot",
			Origin:   0,
			Msgid:    m.MsgID,
			Msgtype:  m.ReplyType,
			SendTime: time.Now(),
			Content:  m.ContentText,
		}
	case *syncmsg.Message: // 同步到的消息
		talk = &schema.Talk{
			OpenKFID: m.OpenKFID,
			ToUserID: m.ExternalUserID,
		}
		msg = &schema.Message{
			MsgFrom:  "sync",
			Origin:   m.Origin,
			Msgid:    m.MsgID,
			Msgtype:  m.MsgType,
			SendTime: time.Unix(int64(m.SendTime), 0),
			Content:  "",
			Raw:      string(m.GetOriginMessage()),
		}
		if m.MsgType == "text" {
			content, _ := m.GetTextMessage()
			msg.Content = content.Text.Content
		}
	default:
		log.Fatalf("save unknown data %v", data)
	}

	talk, err := s.store.Talks().FirstOrCreate(ctx, talk.OpenKFID, talk.ToUserID)
	log.Println("FirstOrCreate err : ", err)
	if err != nil {
		return
	}
	msg.TalkID = talk.ID

	// TODO 保存消息
	err = s.store.Messages().Create(ctx, msg)
	log.Println("s.store.Messages().Create(ctx, msg) err", err)
	if err != nil {
		return
	}

}
