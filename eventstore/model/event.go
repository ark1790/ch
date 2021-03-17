package model

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Event struct {
	ID          string      `json:"id"`
	Email       string      `json:"email"`
	Environment string      `json:"environment"`
	Component   string      `json:"component"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
	CreatedAt   time.Time   `json:"created_at"`
}

func (e *Event) Pre() error {
	if e.ID == "" {
		e.ID = uuid.NewV4().String()
	}

	return nil
}

func (e Event) Marshal() (string, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}
