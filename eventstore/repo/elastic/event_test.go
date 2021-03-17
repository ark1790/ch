package elasticrepo

import (
	"context"
	"testing"
	"time"

	"github.com/ark1790/ch/eventstore/model"
	"github.com/davecgh/go-spew/spew"
	elastic "github.com/olivere/elastic/v7"
	uuid "github.com/satori/go.uuid"
)

func TestCreateEvent(t *testing.T) {
	c, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://localhost:9200"))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = c.Ping("http://localhost:9200").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	m := &model.Event{
		ID:          uuid.NewV4().String(),
		Email:       "m@m2.com",
		Environment: "ENV2",
		Component:   "COMP2",
		Message:     "MESSAGE2",
		Data:        map[string]string{"xx": "yy"},
		CreatedAt:   time.Now(),
	}

	eRepo := NewEventRepo(c)

	if err = eRepo.CreateEvent(context.Background(), m); err != nil {
		t.Fatal(err)
	}

}

func TestListEvents(t *testing.T) {
	c, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://localhost:9200"))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = c.Ping("http://localhost:9200").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	eRepo := NewEventRepo(c)
	qry := map[string]string{
		"environment": "ENV2",
		"from":        "2021-03-17T04:56:28.553159+06:00",
	}
	ets, err := eRepo.FetchEvents(context.Background(), qry, 0, 20)
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(ets)
}
