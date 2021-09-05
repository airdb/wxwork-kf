package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	sdk "github.com/NICEXAI/WeChatCustomerServiceSDK"
	"github.com/NICEXAI/WeChatCustomerServiceSDK/sendmsg"
	"github.com/NICEXAI/WeChatCustomerServiceSDK/syncmsg"
	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/airdb/wxwork-kf/pkg/util"
)

const (
	WelcomeMsg = "您好，这里是宝贝回家公益组织，感谢您的关注和信任。您有寻人、申请志愿者、举报、提供线索、其他咨询等需求，请加宝贝回家唯一全国接待QQ群：1840533。接待群每天9:00-23:00提供咨询登记服务。温馨提示：“宝贝回家”是公益组织，提供的寻亲服务均是免费的，任何发生经济往来的都是假的，请不要相信。"
	DefaultMsg = "回复“帮助”查看更多内容"
)

// Callback - recieve wxkf's notifies.
// @Summary Query item.
// @Description Query item api by id or name.
// @Tags wxkf
// @Accept  json
// @Produce  json
// @Success 200 {string} response "api response"
// @Router /callback [get]
func Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	opts := sdk.CryptoOptions{
		Signature: r.URL.Query().Get("msg_signature"),
		TimeStamp: r.URL.Query().Get("timestamp"),
		Nonce:     r.URL.Query().Get("nonce"),
		EchoStr:   r.URL.Query().Get("echostr"),
	}

	if len(opts.EchoStr) > 0 {
		data, err := app.WxCpt.VerifyURL(opts.Signature, opts.TimeStamp, opts.Nonce, opts.EchoStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)

		return
	}

	body, _ := io.ReadAll(r.Body)
	msg, _ := app.WxCpt.DecryptMsg(opts.Signature, opts.TimeStamp, opts.Nonce, body)

	var msgSyncOpts sdk.SyncMsgOptions
	xml.NewDecoder(bytes.NewReader(msg)).Decode(&msgSyncOpts)
	if cursor, _ := app.Redis.Get(ctx, app.SyncMsgNextCursor).Result(); len(cursor) > 0 {
		msgSyncOpts.Cursor = cursor
	}

	syncMsg, err := app.WxClient.SyncMsg(msgSyncOpts)
	if err == nil {
		app.Redis.Set(ctx, app.SyncMsgNextCursor, syncMsg.NextCursor, 0)
	}

	for _, msg := range syncMsg.MsgList {
		switch msg.Origin {
		case 3: //客户回复的消息
			procUserMsg(ctx, msg)
		case 4: //系统推送的消息
			procSystemMsg(ctx, msg)
		default:
			log.Fatalf("unknown msg origin: %d", msg.Origin)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

// 处理客户回复的消息
func procUserMsg(ctx context.Context, msg syncmsg.Message) {
	var (
		ret          interface{}
		sentCackeKey string // 消息缓存key
		sentCackeTTL time.Duration
	)

	switch msg.MsgType {
	case "text":
		tMsg, _ := msg.GetTextMessage()

		rMsg := &sendmsg.Text{
			Message: sendmsg.Message{
				ToUser:   msg.ExternalUserID,
				OpenKFID: msg.OpenKFID,
				MsgID:    util.RandString(32),
			},
			MsgType: "text",
		}

		if tMsg.Text.Content == "帮助" {
			rMsg.Text.Content = WelcomeMsg
		} else {
			rMsg.Text.Content = DefaultMsg
		}
		ret = rMsg
	default: // 默认回复
		log.Fatalf("unknown user event type: %s", msg.MsgType)
	}

	// 发送消息
	if snedMsg(msg, ret) && len(sentCackeKey) > 0 {
		app.Redis.Set(ctx, sentCackeKey, time.Now().String(), sentCackeTTL).Result()
	}
}

// 处理系统消息
func procSystemMsg(ctx context.Context, msg syncmsg.Message) {
	var (
		ret          interface{}
		sentCackeKey string // 消息缓存key
		sentCackeTTL time.Duration
	)

	switch msg.EventType {
	case "enter_session": // 用户进入会话事件, 一天只发一次
		tMsg, _ := msg.GetEnterSessionEvent()
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
	default:
		log.Fatalf("unknown system event type: %s", msg.EventType)
	}

	// 发送消息
	if snedMsg(msg, ret) && len(sentCackeKey) > 0 {
		app.Redis.Set(ctx, sentCackeKey, time.Now().String(), sentCackeTTL).Result()
	}
}

// 统一发送入口
func snedMsg(msg syncmsg.Message, ret interface{}) bool {
	if info, err := app.WxClient.SendMsg(ret); err == nil {
		params, _ := json.Marshal(ret)
		log.Println("result:", msg.EventType, info.MsgID, ", msg:", string(params))

		return true
	} else {
		log.Println("result:", msg.EventType, ", err:", err.Error())

		return false
	}
}
