package repo

import (
	"context"

	"github.com/ark1790/ch/eventstore/model"
)

type EventRepo interface {
	EnsureIndex() error
	CreateEvent(ctx context.Context, m *model.Event) error
	FetchEvents(ctx context.Context, qry map[string]string, offset, limit int) ([]model.Event, error)
}
