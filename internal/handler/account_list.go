package handler

import (
	"context"
	"net/http"

	sdk "github.com/NICEXAI/WeChatCustomerServiceSDK"
	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/go-chi/render"
)

// KfList - 按场景获取客服列表.
// @Summary Query item.
// @Description Query item api by id or name.
// @Tags wxkf
// @Accept  json
// @Produce  json
// @Success 200 {string} response "api response"
// @Router /account/list [get]
func AccountList(w http.ResponseWriter, r *http.Request) {
	list, err := app.WxClient.AccountList()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, list.BaseModel)
	}

	scene := r.URL.Query().Get("scene")
	if len(scene) == 0 {
		scene = "default"
	}

	retList := make([]map[string]string, 0)
	for _, item := range list.AccountList {
		info, err := app.WxClient.AddContactWay(sdk.AddContactWayOptions{
			OpenKFID: item.OpenKFID,
			Scene:    scene,
		})
		if err != nil {
			continue
		}
		retList = append(retList, map[string]string{
			"name":   item.Name,
			"avatar": item.Avatar,
			"url":    info.URL,
		})
	}

	r = r.WithContext(context.WithValue(r.Context(), render.StatusCtxKey, http.StatusOK))
	render.JSON(w, r, retList)
}
