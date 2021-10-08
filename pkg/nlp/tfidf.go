package nlp

import (
	"errors"
	"fmt"
	"github.com/asmcos/requests"
	"os"
)

func Tfidf(input string) (result string, err error) {
	ak := os.Getenv("TfidfAk")
	url := os.Getenv("TfidfUrl")
	welcomeMsg := os.Getenv("WelcomeMsg")

	req := requests.Requests()
	params := requests.Params{
		"ak":       ak,
		"question": input,
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := req.Get(url, params)
	if err != nil {
		fmt.Println("req err: ", err)
		return
	}
	var json map[string]interface{}
	resp.Json(&json)

	if json["code"].(float64) != 200 {
		err = errors.New("api 错误")
		return
	}
	questionList := json["data"].(map[string]interface{})["Data"].([]interface{})
	if len(questionList) == 0 {
		err = errors.New("api 为空")
		result = welcomeMsg
		return
	}
	result = questionList[0].(map[string]interface{})["answer"].(string)
	return

}
