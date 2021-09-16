package handler

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"

	sdk "github.com/NICEXAI/WeChatCustomerServiceSDK"
	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/airdb/wxwork-kf/internal/service"
	"github.com/airdb/wxwork-kf/internal/store/mysql"
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

	mysqlStore, err := mysql.GetFactoryOr(nil) // TODO
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil)
		return
	}

	replySvc := service.NewReply(mysqlStore)
	for _, msg := range syncMsg.MsgList {
		replySvc.ProcMsg(ctx, msg)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}
