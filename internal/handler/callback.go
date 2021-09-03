package handler

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"

	sdk "github.com/NICEXAI/WeChatCustomerServiceSDK"
	"github.com/NICEXAI/WeChatCustomerServiceSDK/sendmsg"
	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/airdb/wxwork-kf/pkg/util"
)

const (
	WelcomeMsg = "您好，这里是宝贝回家公益组织，感谢您的关注和信任。您有寻人、申请志愿者、举报、提供线索、其他咨询等需求，请加宝贝回家唯一全国接待QQ群：1840533。接待群每天9:00-23:00提供咨询登记服务。温馨提示：“宝贝回家”是公益组织，提供的寻亲服务均是免费的，任何发生经济往来的都是假的，  请不要相信。"
	DefaultMsg = "[开发中]默认消息"
)

// Callback - recieve wxkf's notifies.
// @Summary Query item.
// @Description Query item api by id or name.
// @Tags item
// @Accept  json
// @Produce  json
// @Success 200 {string} response "api response"
// @Router /callback [get]
func Callback(w http.ResponseWriter, r *http.Request) {
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

	syncMsg, _ := app.WxClient.SyncMsg(msgSyncOpts)

	for _, msg := range syncMsg.MsgList {
		// 3-客户回复的消息 4-系统推送的消息
		// if msg.Origin == 3 {
		// 	continue
		// }
		var sMsg interface{}
		switch msg.EventType {
		case "enter_session": // 用户进入会话事件
			eMsg, _ := msg.GetEnterSessionEvent()
			tMsg := &sendmsg.Text{
				Message: sendmsg.Message{
					ToUser:   eMsg.ExternalUserID,
					OpenKFID: eMsg.OpenKFID,
					MsgID:    util.RandString(32),
				},
				MsgType: "text",
			}
			// tMsg.Text.Content = "[开发中]欢迎语"
			tMsg.Text.Content = WelcomeMsg
			sMsg = tMsg
		default: // 默认回复
			tMsg := &sendmsg.Text{
				Message: sendmsg.Message{
					ToUser:   msg.ExternalUserID,
					OpenKFID: msg.OpenKFID,
					MsgID:    util.RandString(32),
				},
				MsgType: "text",
			}
			tMsg.Text.Content = DefaultMsg
			sMsg = tMsg
		}
		if rMsg, err := app.WxClient.SendMsg(sMsg); err == nil {
			params, _ := json.Marshal(sMsg)
			log.Println("result:", msg.EventType, rMsg.MsgID, ", msg:", string(params))
		} else {
			log.Println("result:", msg.EventType, ", err:", err.Error())
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}
