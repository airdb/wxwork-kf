package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	sdk "github.com/NICEXAI/WeChatCustomerServiceSDK"
	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/go-chi/chi/v5"
)

// InviteImageUpload - 上传邀请图片.
// @Summary Query item.
// @Description Query item api by id or name.
// @Tags wxkf
// @Accept  json
// @Produce  json
// @Success 200 {string} response "api response"
// @Router /invite/image/{usedBy:[a-z-]+} [put]
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

	fd, fileHeader, err := r.FormFile("img")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("can not find image: %s", err.Error())))
		return
	}
	fileContent, _ := io.ReadAll(fd)
	base64.StdEncoding.Decode(fileContent, fileContent)

	fileInfo, err := app.WxClient.MediaUpload(sdk.MediaUploadOptions{
		Type:     "image",
		FileName: fileHeader.Filename,
		FileSize: fileHeader.Size,
		File:     bytes.NewReader(fileContent),
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fileBase64 := []byte{}
		base64.StdEncoding.Encode(fileContent, fileBase64)
		w.Write([]byte(fmt.Sprintf(
			"can not upload image: %s, %d, %s, %s",
			fileHeader.Filename, fileHeader.Size, err.Error(), string(fileBase64))))
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
