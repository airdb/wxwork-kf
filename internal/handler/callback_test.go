package handler

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	sdk "github.com/NICEXAI/WeChatCustomerServiceSDK"
	"github.com/NICEXAI/WeChatCustomerServiceSDK/sendmsg"
	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/stretchr/testify/assert"
)

func Test_sendMsg(t *testing.T) {
	file, _ := os.OpenFile("./static/727c801fe83ebe9.jpg", os.O_RDONLY, os.ModePerm)
	fileStat, _ := file.Stat()
	fileInfo, err := app.WxClient.MediaUpload(sdk.MediaUploadOptions{
		Type:     "file",
		FileName: "girl.jpeg",
		FileSize: fileStat.Size(),
		File:     file,
	})
	assert.Nil(t, err)
	log.Println(fileInfo.MediaID)

	sMsg := &sendmsg.Text{
		Message: sendmsg.Message{
			ToUser:   "wm2C5gEQAAex98bnNU1Fm5_U7d18pdVg",
			OpenKFID: "wk2C5gEQAAge69_zMhQgfSor6thQJ8og",
			// MsgID:    "7syeVz1hBT5k39sq9k8M2qHZ48WEfQY5mg5gfZgtKUFF",
			MsgID: "12345678901234567890123456789012",
		},
		MsgType: "text",
	}
	sMsg.Text.Content = "欢迎光临Orzlab"
	// sMsg.Image.MediaID = fileInfo.MediaID

	param, _ := json.Marshal(sMsg)
	log.Println(string(param))
	rMsg, err := app.WxClient.SendMsg(sMsg)
	log.Println("result:", rMsg.MsgID, ", err:", err.Error())
}

func Test_callback(t *testing.T) {
	info, _ := app.WxClient.SyncMsg(sdk.SyncMsgOptions{
		Token: "ENC2rfaK4p7tJXDaeDuGZyqBxzvw4UZSEVHrqQUHrxLobrX",
		Limit: 10,
	})

	for _, msg := range info.MsgList {
		log.Println(msg.EventType)
	}
}

func Test_main(t *testing.T) {
	// list, err := app.WxClient.AccountList()
	// assert.Nil(t, err)
	// assert.NotNil(t, list)

	// infoAdd, err := app.WxClient.AccountAdd(sdk.AccountAddOptions{
	// 	Name: "测试客服",
	// 	// MediaID: "294DpAog3YA5b9rTK4PjjfRfYLO0L5qpDHAJIzhhQ2jAEWjb9i661Q4lk8oFnPtmj",
	// })
	// assert.Nil(t, err)
	// assert.NotNil(t, infoAdd)

	info, err := app.WxClient.AddContactWay(sdk.AddContactWayOptions{
		OpenKFID: "wk2C5gEQAAge69_zMhQgfSor6thQJ8og",
		Scene:    "s-admin",
	})
	assert.Nil(t, err)

	log.Println(info, err)
}
