package mysql

import (
	"fmt"
	"log"
	"sync"

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

func (ds *datastore) Messages() store.MessageStore {
	return newMessages(ds)
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
	once         sync.Once
)

// GetFactoryOr create mysql factory with the given config.
func GetFactoryOr(db *gorm.DB) (store.Factory, error) {
	log.Println("GetFactoryOr")
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
