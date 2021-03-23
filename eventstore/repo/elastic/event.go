package elasticrepo

import (
	"context"
	"encoding/json"

	"github.com/ark1790/ch/eventstore/model"
	"github.com/olivere/elastic/v7"
)

type EventRepo struct {
	er *elastic.Client
}

func NewEventRepo(e *elastic.Client) *EventRepo {
	return &EventRepo{
		er: e,
	}
}

func (e EventRepo) EnsureIndex() error {
	return nil
}

func (e EventRepo) CreateEvent(ctxß context.Context, m *model.Event) error {
	_, err := e.er.Index().
		Index("event").
		Id(m.ID).
		BodyJson(m).
		Do(context.Background())

	return err
}

func (e EventRepo) FetchEvents(ctx context.Context, qry map[string]string, offset, limit int) ([]model.Event, error) {
	eQry := elastic.NewBoolQuery()

	rQry := elastic.NewRangeQuery("created_at")
	for k, v := range qry {
		if k == "from" {
			rQry.Gte(v)
		} else if k == "to" {
			rQry.Lte(v)
		} else {
			eQry = eQry.Must(elastic.NewMatchPhraseQuery(k, v))
		}
	}

	eQry = eQry.Must(rQry)

	sResult, err := e.er.Search().Index("event").
		Query(eQry).
		From(offset).
		Size(limit).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	ets := []model.Event{}
	evt := model.Event{}

	for _, hit := range sResult.Hits.Hits {
		err := json.Unmarshal(hit.Source, &evt)
		if err != nil {
			return nil, err
		}

		ets = append(ets, evt)
	}

	return ets, nil
}
