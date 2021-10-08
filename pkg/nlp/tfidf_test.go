package nlp_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/airdb/wxwork-kf/pkg/nlp"
)

func TestTfidf(t *testing.T) {
	os.Setenv("TfidfAk", "de02e765922b8d7bc2c0d1b0d2dfe358")
	os.Setenv("TfidfUrl", "http://127.0.0.1:8000/v1/api/question/match")
	os.Setenv("WelcomeMsg", "您好，这里是宝贝回家公益组织，感谢您的关注和信任。您有寻人、申请志愿者、举报、提供线索、其他咨询等需求，请加宝贝回家唯一全国接待QQ群：1840533。接待群每天9:00-23:00提供咨询登记服务。温馨提示：“宝贝回家”是公益组织，提供的寻亲服务均是免费的，任何发生经济往来的都是假的，  请不要相信。")
	fmt.Println(nlp.Tfidf("教师怎么离职"))
}
