package store

import (
	"context"

	"github.com/airdb/wxwork-kf/pkg/schema"
)

// TalkStore defines the talk storage interface.
type TalkStore interface {
	Create(ctx context.Context, talk *schema.Talk) error
}
