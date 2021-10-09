package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/airdb/sailor/dbutil"
	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/airdb/wxwork-kf/internal/service"
	"github.com/airdb/wxwork-kf/internal/store/mysql"
	"github.com/silenceper/wechat/v2/work/kf"
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
	opts := kf.SignatureOptions{
		Signature: r.URL.Query().Get("msg_signature"),
		TimeStamp: r.URL.Query().Get("timestamp"),
		Nonce:     r.URL.Query().Get("nonce"),
		EchoStr:   r.URL.Query().Get("echostr"),
	}

	if len(opts.EchoStr) > 0 {
		data, err := app.WxWorkKF.VerifyURL(opts)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Panicf("can not read request body: %s", err.Error())
	}
	log.Println("callback body: ", string(body))
	cbMsg, err := app.WxWorkKF.GetCallbackMessage(body)
	if err != nil {
		log.Panicf("can not decrypt callback\n")
	}

	syncMsgOpts := kf.SyncMsgOptions{Token: cbMsg.Token}
	// 获取上次消息游标
	cursor, err := app.Redis.Get(ctx, app.SyncMsgNextCursor).Result()
	if err == nil && len(cursor) > 0 {
		syncMsgOpts.Cursor = cursor
	}

	app.Redis.Set(ctx, app.SyncMsgNextCursor, cbMsg.Token, 0)
	syncMsg, err := app.WxWorkKF.SyncMsg(syncMsgOpts)
	if err == nil {
		// 保存本次消息游标
		app.Redis.Set(ctx, app.SyncMsgNextCursor, syncMsg.NextCursor, 0)
	}

	mysqlStore, err := mysql.GetFactoryOr(dbutil.WriteDefaultDB()) // TODO
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil)
		return
	}

	replySvc := service.NewReply(mysqlStore)
	log.Println("replySvc",replySvc)
	for _, msg := range syncMsg.MsgList {
		log.Println("sync from wechat, msg:", msg)
		replySvc.ProcMsg(ctx, msg)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}
