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

// FirstOrCreate find first talk by kfid and userid.
func (u *talks) FirstOrCreate(ctx context.Context, openKFID, toUserID string) (*schema.Talk, error) {
	talk := schema.Talk{
		OpenKFID: openKFID,
		ToUserID: toUserID,
	}
	if err := u.db.Where(&talk).First(&talk).Error; err == nil {
		return &talk, nil
	}
	if len(talk.ID) >0 {
		return &talk, nil
	}
	if err := u.db.Create(&talk).Error; err == nil {
		return &talk, nil
	} else {
		return nil, err
	}
}
