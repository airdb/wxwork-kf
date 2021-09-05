package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	sdk "github.com/NICEXAI/WeChatCustomerServiceSDK"
	"github.com/NICEXAI/WeChatCustomerServiceSDK/cache"
	"github.com/NICEXAI/WeChatCustomerServiceSDK/crypto"
	"github.com/NICEXAI/WeChatCustomerServiceSDK/sendmsg"
	"github.com/airdb/sailor/deployutil"
	"github.com/airdb/sailor/faas"
	"github.com/airdb/sailor/version"

	// "github.com/airdb/wxwork-kf/pkg/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs" // docs is generated by Swag CLI, you have to import it.
)

var (
	err      error
	wxOption sdk.Options
	wxCpt    *crypto.WXBizMsgCrypt
	wxClient *sdk.Client
)

func init() {
	wxOption.CorpID = os.Getenv("WXKF_CORP_ID")
	wxOption.Secret = os.Getenv("WXKF_SECRET")
	wxOption.Token = os.Getenv("WXKF_SECRET")
	wxOption.EncodingAESKey = os.Getenv("WXKF_ENCODING_AES_KEY")

	var wxRedisOption cache.RedisOptions
	if redisDB, err := strconv.Atoi(os.Getenv("WXKF_REDIS_DB")); err == nil {
		wxRedisOption.Addr = os.Getenv("WXKF_REDIS_ADDR")
		wxRedisOption.Password = os.Getenv("WXKF_REDIS_PASSWD")
		wxRedisOption.DB = redisDB
		wxOption.Cache = cache.NewRedis(wxRedisOption)
	} else {
		log.Panicf("init redis err: %s", err.Error())
	}

	log.Printf("get env from serverless.yaml: redis: %s,  corpID: %s, secret: %s, token: %s, encodingAESKey: %s\n",
		wxRedisOption.Addr,
		wxOption.CorpID,
		wxOption.Secret,
		wxOption.Token,
		wxOption.EncodingAESKey,
	)

	wxCpt = crypto.NewWXBizMsgCrypt(
		wxOption.Token, wxOption.EncodingAESKey, wxOption.CorpID, crypto.XmlType,
	)
	if wxClient, err = sdk.New(wxOption); err != nil {
		log.Panicf("init client err: %s", err.Error())
	}
}

// @title Airdb Serverlesss Example API
// @version 0.0.1
// @description This is a sample server Petstore server.
// @termsOfService https://airdb.github.io/terms.html

// @contact.name API Support
// @contact.url https://work.weixin.qq.com/kfid/kfc5fdb2e0a1f297753
// @contact.email info@airdb.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @info.x-logo.url: http://www.apache.org/licenses/LICENSE-2.0.html

// @host service-iw6drlfr-1251018873.sh.apigw.tencentcs.com
// @BasePath /wxkf
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeHTML))

	r.Route("/wxkf", func(r chi.Router) {
		r.Get("/version", faas.HandleVersion)
		r.HandleFunc("/callback", HandleCallback)
	})

	fmt.Println("hello", deployutil.GetDeployStage())

	faas.RunTencentChiWithSwagger(r)
}

const (
	WelcomeMsg = "您好，这里是宝贝回家公益组织，感谢您的关注和信任。您有寻人、申请志愿者、举报、提供线索、其他咨询等需求，请加宝贝回家唯一全国接待QQ群：1840533。接待群每天9:00-23:00提供咨询登记服务。温馨提示：“宝贝回家”是公益组织，提供的寻亲服务均是免费的，任何发生经济往来的都是假的，  请不要相信。"
	DefaultMsg = "[开发中]默认消息"
)

// HandleCallback - recieve wxkf's notifies.
// @Summary Query item.
// @Description Query item api by id or name.
// @Tags item
// @Accept  json
// @Produce  json
// @Success 200 {string} response "api response"
// @Router /callback [get]
func HandleCallback(w http.ResponseWriter, r *http.Request) {
	opts := sdk.CryptoOptions{
		Signature: r.URL.Query().Get("msg_signature"),
		TimeStamp: r.URL.Query().Get("timestamp"),
		Nonce:     r.URL.Query().Get("nonce"),
		EchoStr:   r.URL.Query().Get("echostr"),
	}

	if len(opts.EchoStr) > 0 {
		data, err := wxCpt.VerifyURL(opts.Signature, opts.TimeStamp, opts.Nonce, opts.EchoStr)
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
	msg, _ := wxCpt.DecryptMsg(opts.Signature, opts.TimeStamp, opts.Nonce, body)

	var msgSyncOpts sdk.SyncMsgOptions
	xml.NewDecoder(bytes.NewReader(msg)).Decode(&msgSyncOpts)

	syncMsg, _ := wxClient.SyncMsg(msgSyncOpts)

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
					MsgID:    RandString(32),
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
					MsgID:    RandString(32),
				},
				MsgType: "text",
			}
			tMsg.Text.Content = DefaultMsg
			sMsg = tMsg
		}
		if rMsg, err := wxClient.SendMsg(sMsg); err == nil {
			params, _ := json.Marshal(sMsg)
			log.Println("result:", msg.EventType, rMsg.MsgID, ", msg:", string(params))
		} else {
			log.Println("result:", msg.EventType, ", err:", err.Error())
		}
	}

	w.WriteHeader(http.StatusOK)
}

func myMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(version.GetBuildInfo().ToString()))
		// Here we are pssing our custom response writer to the next http handler.
		next.ServeHTTP(w, r)

		// Here we are adding our custom stuff to the response, which we received after http handler execution.
		// myResponseWriter.buf.WriteString(" and some additional modifications")
	})
}

// func RandString(len int) string {
// 	r:=rand.New(rand.NewSource(time.Now().Unix()))
// 	bytes := make([]byte, len)
// 	for i := 0; i < len; i++ {
// 			b := r.Intn(26) + 65
// 			bytes[i] = byte(b)
// 	}
// 	return string(bytes)
// }

func RandString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var src = rand.NewSource(time.Now().UnixNano())

	const (
		letterIdxBits = 6
		letterIdxMask = 1<<letterIdxBits - 1
		letterIdxMax  = 63 / letterIdxBits
	)
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
