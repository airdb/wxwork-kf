package mysql_test

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/airdb/wxwork-kf/internal/store/mysql"
)

var (
	once sync.Once
)

func TestGetConnection(t *testing.T)  {
	dns := "demo:demo@123@tcp(sh-cdb-n7qw1jqg.sql.tencentcdb.com:62974)/demo"
	os.Setenv("MAIN_DSN_WRITE",dns)
	os.Setenv("MAIN_DSN_READ",dns)

	fmt.Println(os.Getenv("MAIN_DSN_WRITE"))
	mysql.GetConnection()
}

func a()  {
	fmt.Println("a")
}

func b()  {
	fmt.Println("b")
}
func TestOnce(t *testing.T){
	once.Do(func() {
		a()
	})
	once.Do(func() {
		b()
	})
}
