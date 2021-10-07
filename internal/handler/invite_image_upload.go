package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/silenceper/wechat/v2/officialaccount/material"
)

// InviteImageUpload - 上传邀请图片.
// @Summary Query item.
// @Description Query item api by id or name.
// @Tags wxkf
// @Accept  json
// @Produce  json
// @Success 200 {string} response "api response"
// @Router /invite/image/{usedBy:} [put]
func InviteImageUpload(w http.ResponseWriter, r *http.Request) {
	var err error

	err = r.ParseMultipartForm(2 << 20) // 8 MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can not parse form"))
		return
	}

	usedBy := chi.URLParam(r, "usedBy")
	if len(usedBy) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can not find url param"))
		return
	}
	inviteImageCacheKey := strings.Join([]string{
		app.InviteImagePrefix, usedBy,
	}, ":")

	// 解 base64 后保存在临时文件中上传
	fd, _, err := r.FormFile("img")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("can not find image: %s", err.Error())))
		return
	}
	fileBase64, _ := io.ReadAll(fd)
	fileContent := make([]byte, base64.StdEncoding.DecodedLen(len(fileBase64)))
	n, err := base64.StdEncoding.Decode(fileContent, fileBase64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("can not decode image(%d): %s", n, err.Error())))
		return
	}
	tmpFile, _ := ioutil.TempFile("", "tmp")
	io.Copy(tmpFile, bytes.NewBuffer(fileContent))
	tmpFile.Sync()
	tmpStat, _ := tmpFile.Stat()

	fileInfo, err := app.WxWorkMedia.MediaUpload(material.MediaTypeImage, tmpFile.Name())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(
			"can not upload image: %s, %d, %s",
			tmpFile.Name(), tmpStat.Size(), err.Error())))
		return
	}

	_, err = app.Redis.Set(r.Context(), inviteImageCacheKey, fileInfo.MediaID, 0).Result()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("can not cache media: %s", err.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fileInfo.MediaID))
}
