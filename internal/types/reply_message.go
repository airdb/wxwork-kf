package types

import (
	"encoding/json"

	sdk "github.com/NICEXAI/WeChatCustomerServiceSDK"
	"github.com/NICEXAI/WeChatCustomerServiceSDK/sendmsg"
)

type ContentMenu struct {
	HeadContent string
	List        []interface{}
	TailContent string
}

type ReplyMessage struct {
	ToUser   string
	OpenKFID string
	MsgID    string

	ReplyType    string
	ContentText  string
	ContentImage string
	ContentMenu  ContentMenu

	ActionTransState    int
	ActionTransServicer string

	msg interface{} // 组装后的消息体
}

func NewReplyMessage(toUser, openKFID, msgID string) *ReplyMessage {
	return &ReplyMessage{
		ToUser: toUser, OpenKFID: openKFID, MsgID: msgID,
	}
}

// Assume 用于组成微信客服接口请求体
func (m ReplyMessage) Assume() (interface{}, error) {
	return m.msg, nil
}

func (m *ReplyMessage) SetText(s string) {
	m.ReplyType = WxMsgTypeText
	m.ContentText = s

	msg := sendmsg.Text{
		Message: m.getMessageHead(),
		MsgType: m.ReplyType,
	}
	msg.Text.Content = m.ContentText
	m.msg = msg
}

func (m *ReplyMessage) SetImage(s string) {
	m.ReplyType = WxMsgTypeImg
	m.ContentImage = s

	msg := sendmsg.Image{
		Message: m.getMessageHead(),
		MsgType: m.ReplyType,
	}
	msg.Image.MediaID = m.ContentImage
	m.msg = msg
}

func (m *ReplyMessage) SetMenu(cm ContentMenu) {
	m.ReplyType = WxMsgTypeMenu
	m.ContentMenu = cm

	msg := sendmsg.Menu{
		Message: m.getMessageHead(),
		MsgType: m.ReplyType,
	}
	msg.MsgMenu.HeadContent = m.ContentMenu.HeadContent
	msg.MsgMenu.List = m.ContentMenu.List
	msg.MsgMenu.TailContent = m.ContentMenu.TailContent
	m.msg = msg
}

func (m *ReplyMessage) SetActionTrans(state int, servicer string) {
	m.ReplyType = WxMsgTypeActionTrans
	m.ActionTransState = state
	m.ActionTransServicer = servicer

	msg := sdk.ServiceStateTransOptions{
		OpenKFID:       m.OpenKFID,
		ExternalUserID: m.ToUser,
		ServiceState:   m.ActionTransState,
		ServicerUserID: m.ActionTransServicer,
	}
	m.msg = msg
}

// Content 需要记录消息内容
func (m ReplyMessage) Content() string {
	var content string
	switch m.ReplyType {
	case WxMsgTypeText:
		content = m.ContentText
	case WxMsgTypeImg:
		content = m.ContentImage
	case WxMsgTypeVideo:
		fallthrough
	case WxMsgTypeVoice:
		fallthrough
	case WxMsgTypeFile:
		fallthrough
	case WxMsgTypeLocation:
		content = ""
	case WxMsgTypeMenu:
		if bs, err := json.Marshal(m.ContentMenu); err == nil {
			content = string(bs)
		} else {
			content = ""
		}
	case WxMsgTypeActionTrans:
		if bs, err := json.Marshal(m.ContentMenu); err == nil {
			content = string(bs)
		} else {
			content = ""
		}
	default:
		content = ""
	}

	return content
}

func (m ReplyMessage) getMessageHead() sendmsg.Message {
	return sendmsg.Message{
		ToUser:   m.ToUser,
		OpenKFID: m.OpenKFID,
		MsgID:    m.MsgID,
	}
}
