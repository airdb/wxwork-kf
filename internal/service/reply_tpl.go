package service

import (
	"log"
	"strings"

	"github.com/NICEXAI/WeChatCustomerServiceSDK/syncmsg"
	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/airdb/wxwork-kf/internal/types"
	"github.com/airdb/wxwork-kf/pkg/util"
)

const (
	WelcomeMsg = "您好，这里是宝贝回家公益组织，感谢您的关注和信任。您有寻人、申请志愿者、举报、提供线索、其他咨询等需求，请加宝贝回家唯一全国接待QQ群：1840533。接待群每天9:00-23:00提供咨询登记服务。温馨提示：“宝贝回家”是公益组织，提供的寻亲服务均是免费的，任何发生经济往来的都是假的，请不要相信。"
	DefaultMsg = "回复“帮助”查看更多内容"
)

// 消息匹配方式
type MatchMethod int

const (
	MatchMethodFull     MatchMethod = iota // 匹配全文
	MatchMethodKeyword                     // 匹配关键字
	MatchMethodPrefix                      // 匹配前缀
	MatchMethodRegexp                      // 正则匹配
	MatchMethodImg                         // 匹配图片
	MatchMethodVideo                       // 匹配视频
	MatchMethodVoice                       // 匹配音频
	MatchMethodFile                        // 匹配文件
	MatchMethodLocation                    // 匹配定位
)

// ReplyType 返回消息类型
type ReplyType int

const (
	ReplyTypeText        ReplyType = iota       // 返回文本消息
	ReplyTypeImage                              // 返回图片消息
	ReplyTypeMenu                               // 返回菜单消息
	ReplyTypeActionTrans ReplyType = iota + 100 // 分配客服会话
)

// 填充式消息
type ReplyCallback func(mr *types.ReplayMessage) error

var (
	// TplReplys 根据用户消息内容，返回对话内容
	TplReplys = &ReplyTpls{
		{"text", MatchMethodFull, "default", ReplyTypeText, WelcomeMsg},
		{"text", MatchMethodFull, "帮助", ReplyTypeMenu, ReplyCallback(func(mr *types.ReplayMessage) error {
			cm := types.ContentMenu{
				HeadContent: WelcomeMsg,
				List: []interface{}{
					map[string]interface{}{
						"type":  "click",
						"click": map[string]string{"id": "welcome", "content": "欢迎语"},
					},
					map[string]interface{}{
						"type":  "click",
						"click": map[string]string{"id": "search", "content": "寻人"},
					},
					map[string]interface{}{
						"type":  "click",
						"click": map[string]string{"id": "clue", "content": "线索"},
					},
					map[string]interface{}{
						"type":  "click",
						"click": map[string]string{"id": "volunteer", "content": "志愿者"},
					},
				},
			}
			mr.SetMenu(cm)
			return nil
		})},
		{"text", MatchMethodFull, "寻人", ReplyTypeText, "[寻人](在这里输入你能提供的信息)"},
		{"text", MatchMethodFull, "线索", ReplyTypeText, "[线索](在这里输入你能提供的线索)"},
		{"text", MatchMethodFull, "人工客服", ReplyTypeActionTrans, ""},
		{"text", MatchMethodPrefix, "[线索]", ReplyTypeActionTrans, ""},
		{"text", MatchMethodPrefix, "[志愿者]", ReplyTypeImage, "图片ID"},
		// TODO
		{"image", MatchMethodImg, "【图片消息】", ReplyTypeText, "【图片消息】"},
		{"video", MatchMethodVideo, "【视频消息】", ReplyTypeText, "【视频消息】"},
		{"voice", MatchMethodVoice, "【语音消息】", ReplyTypeText, "【语音消息】"},
		{"file", MatchMethodFile, "【文件消息】", ReplyTypeText, "【文件消息】"},
		{"location", MatchMethodLocation, "【位置消息】", ReplyTypeText, "【位置消息】"},
	}
)

type ReplyTpl struct {
	MatchType   string      // 消息配置类型, 可选值 text, image, video, voice, file, location
	MatchMethod MatchMethod // 消息配置方式, 可选值: full, keyword, prefix
	MatchValue  string      // 消息配置内容

	ReplyType ReplyType   // 返回消息类型, 可选值: action, text
	Message   interface{} // 消息内容 or interface
}

// Gen 组装消息
func (rt ReplyTpl) Gen(toUser, openKFID string) *types.ReplayMessage {
	ret := types.NewReplayMessage(toUser, openKFID, util.RandString(32))

	switch rt.ReplyType {
	case ReplyTypeText: // 文本消息
		ret.SetText(rt.Message.(string))
	case ReplyTypeImage: // 图片消息
		ret.SetImage(rt.Message.(string))
	case ReplyTypeMenu: // 菜单消息
		callback, ok := rt.Message.(ReplyCallback)
		if ok {
			callback(ret)
		}
	case ReplyTypeActionTrans:
		// 查找客服账号列表
		accountList, err := app.WxClient.AccountList()
		if err != nil || len(accountList.AccountList) == 0 {
			return nil
		}
		account := accountList.AccountList[0]

		// 接待人员列表
		receptionisList, err := app.WxClient.ReceptionistList(account.OpenKFID)
		if err != nil || len(receptionisList.ReceptionistList) == 0 {
			return nil
		}
		receptionis := receptionisList.ReceptionistList[0]

		ret.SetActionTrans(3, receptionis.UserID)
	}

	return ret
}

// Match 根据不同的消息类型选择不同的匹配方式
func (rt ReplyTpl) Match(msg syncmsg.Message) bool {
	switch msg.MsgType {
	case types.WxMsgTypeText:
		if info, err := msg.GetTextMessage(); err == nil {
			return rt.matchText(info.Text.Content)
		}
		return false
	case types.WxMsgTypeImg: // 图片
		return true
	case types.WxMsgTypeVideo: // 视频
		return true
	case types.WxMsgTypeVoice: // 语音
		return true
	case types.WxMsgTypeFile: // 文件
		return true
	case types.WxMsgTypeLocation: // 位置
		return true
	default: // 默认回复
		log.Fatalf("unknown user msg type: %s", msg.MsgType)
		return false
	}
}

func (rt ReplyTpl) matchText(s string) bool {
	switch rt.MatchMethod {
	case MatchMethodFull:
		return rt.MatchValue == s
	case MatchMethodKeyword:
		return strings.Contains(s, rt.MatchValue)
	case MatchMethodPrefix:
		return strings.HasPrefix(s, rt.MatchValue)
	default:
		return false
	}
}

type ReplyTpls []*ReplyTpl

// Match 查询用户消息是否命中
func (rts ReplyTpls) Match(msg syncmsg.Message) (*ReplyTpl, bool) {
	for _, rt := range rts {
		if rt.Match(msg) {
			return rt, true
		}
	}
	return nil, false
}

func (rts ReplyTpls) MatchText(msg string) (*ReplyTpl, bool) {
	for _, rt := range rts {
		if rt.matchText(msg) {
			return rt, true
		}
	}
	return nil, false
}

// Default 默认消息
func (rts ReplyTpls) Default() (*ReplyTpl, bool) {
	for _, rt := range rts {
		if rt.matchText("default") {
			return rt, true
		}
	}
	return nil, false
}
