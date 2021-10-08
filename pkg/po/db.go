package po

import (
	"github.com/airdb/sailor/dbutil"
)

// InitDB  初始化db
func InitDB() {
	dbutil.InitDefaultDB()
}
