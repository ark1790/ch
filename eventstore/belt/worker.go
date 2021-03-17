package belt

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ark1790/ch/eventstore/model"
	"github.com/ark1790/ch/eventstore/repo"
)

type Worker struct {
	eventRepo repo.EventRepo
	queue     Queue
}

func NewWorker(q Queue, er repo.EventRepo) *Worker {
	return &Worker{
		queue:     q,
		eventRepo: er,
	}
}

func (w *Worker) Start() {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.Start()
		}
	}()

	for {
		func() {

			ack := make(chan bool)
			msg, err := w.queue.Pop(ack)
			if msg == nil && err == nil {
				close(ack)
				return
			} else if err != nil {
				log.Println(err)
			}

			evt := &model.Event{}
			if err := json.Unmarshal([]byte(*msg), evt); err != nil {
				log.Println(err)
				return
			}

			log.Println("Got Message for Event-ID:", evt.ID)

			if err := w.eventRepo.CreateEvent(context.Background(), evt); err != nil {
				log.Println(err)
				return
			}

			log.Println("Created Event For Event-ID:", evt.ID)

		}()
	}
}
