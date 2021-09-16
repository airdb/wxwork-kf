package mysql

import (
	"context"

	"github.com/airdb/wxwork-kf/pkg/schema"
	"gorm.io/gorm"
)

type talks struct {
	db *gorm.DB
}

func newTalks(ds *datastore) *talks {
	return &talks{db: ds.db}
}

// Create creates a new talk item.
func (u *talks) Create(ctx context.Context, talk *schema.Talk) error {
	return u.db.Create(&talk).Error
}
