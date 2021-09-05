package mysql

import (
	"context"

	"github.com/airdb/wxwork-kf/pkg/schema"
	"gorm.io/gorm"
)

type talks struct {
	writeDb *gorm.DB
	readDb *gorm.DB
}

func newTalks(ds *datastore) *talks {
	return &talks{
		writeDb: ds.writeDb,
		readDb: ds.readDb,
	}
}

// Create creates a new talk item.
func (u *talks) Create(ctx context.Context, talk *schema.Talk) error {
	return u.writeDb.Create(&talk).Error
}
