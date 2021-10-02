package util

import (
	"encoding/json"
	"fmt"

	"github.com/silenceper/wechat/v2/officialaccount/context"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"github.com/silenceper/wechat/v2/util"
)

const (
	mediaUploadURL      = "https://qyapi.weixin.qq.com/cgi-bin/media/upload"
	mediaUploadImageURL = "https://qyapi.weixin.qq.com/cgi-bin/media/uploadimg"
	mediaGetURL         = "https://qyapi.weixin.qq.com/cgi-bin/media/get"
)

// Material 素材管理
type Material struct {
	*context.Context
}

// NewMaterial init
func NewMaterial(context *context.Context) *Material {
	material := new(Material)
	material.Context = context
	return material
}

// Media 临时素材上传返回信息
type Media struct {
	util.CommonError

	Type      material.MediaType `json:"type"`
	MediaID   string             `json:"media_id"`
	CreatedAt string             `json:"created_at"`
}

// MediaUpload 临时素材上传
func (material *Material) MediaUpload(mediaType material.MediaType, filename string) (media Media, err error) {
	var accessToken string
	accessToken, err = material.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s&type=%s", mediaUploadURL, accessToken, mediaType)
	var response []byte
	response, err = util.PostFile("media", filename, uri)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &media)
	if err != nil {
		return
	}
	if media.ErrCode != 0 {
		err = fmt.Errorf("MediaUpload error : errcode=%v , errmsg=%v", media.ErrCode, media.ErrMsg)
		return
	}
	return
}

// GetMediaURL 返回临时素材的下载地址供用户自己处理
// NOTICE: URL 不可公开，因为含access_token 需要立即另存文件
func (material *Material) GetMediaURL(mediaID string) (mediaURL string, err error) {
	var accessToken string
	accessToken, err = material.GetAccessToken()
	if err != nil {
		return
	}
	mediaURL = fmt.Sprintf("%s?access_token=%s&media_id=%s", mediaGetURL, accessToken, mediaID)
	return
}

// resMediaImage 图片上传返回结果
type resMediaImage struct {
	util.CommonError

	URL string `json:"url"`
}

// ImageUpload 图片上传
func (material *Material) ImageUpload(filename string) (url string, err error) {
	var accessToken string
	accessToken, err = material.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s", mediaUploadImageURL, accessToken)
	var response []byte
	response, err = util.PostFile("media", filename, uri)
	if err != nil {
		return
	}
	var image resMediaImage
	err = json.Unmarshal(response, &image)
	if err != nil {
		return
	}
	if image.ErrCode != 0 {
		err = fmt.Errorf("UploadImage error : errcode=%v , errmsg=%v", image.ErrCode, image.ErrMsg)
		return
	}
	url = image.URL
	return
}
