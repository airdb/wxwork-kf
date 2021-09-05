package mysql

import (
	"fmt"
	"github.com/airdb/sailor/dbutil"
	"github.com/airdb/wxwork-kf/pkg/po"
	"sync"

	"github.com/airdb/wxwork-kf/internal/store"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type datastore struct {
	writeDb *gorm.DB
	readDb *gorm.DB
}

func (ds *datastore) Talks() store.TalkStore {
	return newTalks(ds)
}

func (ds *datastore) Close() error {
	writeDb, wdbErr := ds.writeDb.DB()
	if wdbErr != nil {
		return errors.Wrap(wdbErr, "get gorm writeDB instance failed")
	}

	readDb, rdbErr := ds.readDb.DB()
	if rdbErr != nil {
		return errors.Wrap(rdbErr, "get gorm readDB instance failed")
	}

	wdbErr = writeDb.Close()
	if wdbErr != nil {
		return wdbErr
	}
	rdbErr = readDb.Close()
	if rdbErr != nil {
		return rdbErr
	}
	return nil
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

// GetFactoryOr create mysql factory with the given config.
func GetFactoryOr() (store.Factory, error) {
	var err error
	once.Do(func() {
		po.InitDB()
		mysqlFactory = &datastore{
			writeDb:dbutil.WriteDefaultDB(),
			readDb: dbutil.ReadDefaultDB(),
		}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %s", mysqlFactory, err.Error())
	}

	return mysqlFactory, nil
}
