package handler

import (
	"context"
	"net/http"

	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/go-chi/render"
	"github.com/silenceper/wechat/v2/work/kf"
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
	list, err := app.WxWorkKF.AccountList()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, list.CommonError)
	}

	scene := r.URL.Query().Get("scene")
	if len(scene) == 0 {
		scene = "default"
	}

	retList := make([]map[string]string, 0)
	for _, item := range list.AccountList {
		info, err := app.WxWorkKF.AddContactWay(kf.AddContactWayOptions{
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
