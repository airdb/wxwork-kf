package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// EchoFile - 上传文件测试，上传图片时需要base64，原因未知.
// @Summary Query item.
// @Description Query item api by id or name.
// @Tags wxkf
// @Accept  json
// @Produce  json
// @Success 200 {string} response "api response"
// @Router /echo/file [put]
func EchoFile(w http.ResponseWriter, r *http.Request) {
	var err error

	file, _, err := r.FormFile("img")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("can not find image: %s", err.Error())))
		return
	}

	fileBytes, _ := ioutil.ReadAll(file)

	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
}
