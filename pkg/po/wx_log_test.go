package po_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/airdb/wxwork-kf/pkg/po"
)

func TestWxKfLogSave(t *testing.T) {
	os.Setenv("MAIN_DSN_WRITE", "demo:demo@123@tcp(sh-cdb-n7qw1jqg.sql.tencentcdb.com:62974)/demo")
	os.Setenv("MAIN_DSN_READ", "demo:demo@123@tcp(sh-cdb-n7qw1jqg.sql.tencentcdb.com:62974)/demo")
	po.InitDB()
	data := new(po.WxKfLog)
	data.Input = "test"
	data.Response = "test"
	data.OpenKFID = "test"
	data.ToUserID = "test"
	data.MsgID = "test"
	err := po.WxKfLogSave(data)
	fmt.Println(err)
}
