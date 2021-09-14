package types

import (
	sdk "github.com/NICEXAI/WeChatCustomerServiceSDK"
	"github.com/NICEXAI/WeChatCustomerServiceSDK/sendmsg"
)

type ContentMenu struct {
	HeadContent string
	List        []interface{}
	TailContent string
}

type ReplayMessage struct {
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

func NewReplayMessage(toUser, openKFID, msgID string) *ReplayMessage {
	return &ReplayMessage{
		ToUser: toUser, OpenKFID: openKFID, MsgID: msgID,
	}
}

// Assume 用于组成微信客服接口请求体
func (m ReplayMessage) Assume() (interface{}, error) {
	return m.msg, nil
}

func (m *ReplayMessage) SetText(s string) {
	m.ReplyType = "text"
	m.ContentText = s

	msg := sendmsg.Text{
		Message: m.getMessageHead(),
		MsgType: m.ReplyType,
	}
	msg.Text.Content = m.ContentText
	m.msg = msg
}

func (m *ReplayMessage) SetImage(s string) {
	m.ReplyType = "image"
	m.ContentImage = s

	msg := sendmsg.Image{
		Message: m.getMessageHead(),
		MsgType: m.ReplyType,
	}
	msg.Image.MediaID = m.ContentImage
	m.msg = msg
}

func (m *ReplayMessage) SetMenu(cm ContentMenu) {
	m.ReplyType = "msgmenu"
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

func (m *ReplayMessage) SetActionTrans(state int, servicer string) {
	m.ReplyType = "actionTrans"
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

func (m ReplayMessage) getMessageHead() sendmsg.Message {
	return sendmsg.Message{
		ToUser:   m.ToUser,
		OpenKFID: m.OpenKFID,
		MsgID:    m.MsgID,
	}
}
