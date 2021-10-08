package mysql

import (
	"fmt"
	"log"
	"sync"

	"github.com/airdb/sailor/dbutil"
	"github.com/airdb/wxwork-kf/internal/store"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Talks() store.TalkStore {
	return newTalks(ds)
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

var (
	mysqlFactory store.Factory
	mysqlOnce         sync.Once
	once         sync.Once
	readDB		*gorm.DB
	writeDB		*gorm.DB
)

// GetFactoryOr create mysql factory with the given config.
func GetFactoryOr(db *gorm.DB) (store.Factory, error) {
	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		dbIns = db
		mysqlFactory = &datastore{dbIns}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %s", mysqlFactory, err.Error())
	}

	return mysqlFactory, nil
}

// GetConnection  get mysql connection, default is write DB
func GetConnection()*gorm.DB{
	mysqlOnce.Do(func() {
		dbutil.InitDefaultDB()
		writeDB =  dbutil.WriteDefaultDB()
		readDB =  dbutil.ReadDefaultDB()
	})
	if writeDB == nil {
		log.Println("mysql 连接失败")
		panic("mysql 连接失败")
	}
	return writeDB
}
