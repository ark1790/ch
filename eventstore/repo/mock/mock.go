package mockrepo

import (
	"context"

	"github.com/ark1790/ch/eventstore/model"
)

type MockRepo struct {
	Events []model.Event
}

func NewEventRepo() *MockRepo {
	return &MockRepo{
		Events: []model.Event{},
	}
}

func (mr *MockRepo) EnsureIndex() error {
	return nil
}

func (mr *MockRepo) CreateEvent(ctx context.Context, m *model.Event) error {
	mr.Events = append(mr.Events, *m)
	return nil
}

func (mr MockRepo) FetchEvents(ctx context.Context, qry map[string]string, offset, limit int) ([]model.Event, error) {
	return mr.Events, nil
}
