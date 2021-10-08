package mysql

import (
	"context"
	"log"

	"github.com/airdb/wxwork-kf/pkg/schema"
	"gorm.io/gorm"
)

type messages struct {
	db *gorm.DB
}

func newMessages(ds *datastore) *messages {
	return &messages{db: ds.db}
}

// Create creates a new message item.
func (m *messages) Create(ctx context.Context, message *schema.Message) error {
	log.Println("record message, value: ", message)
	return m.db.Create(&message).Error
}
