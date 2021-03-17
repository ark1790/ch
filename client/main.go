package main

import (
	"context"
	"log"

	"time"

	"github.com/ark1790/ch/eventstore/proto"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

func main() {
	conn, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewEventStoreClient(conn)

	data := map[string]interface{}{
		"A": "X",
		"B": map[string]interface{}{
			"C": "D",
		},
	}

	d, _ := structpb.NewStruct(data)
	req := &proto.ReqCreateEvent{
		Email:       "ABC@B.com",
		Environment: "Aenvironment",
		Component:   "A COMPONENT",
		Message:     "the buyer #123456 has placed an order successfully",
		Data:        d,
		CreatedAt:   time.Now().Unix(),
	}

	resp, err := client.CreateEvent(context.Background(), req)
	if err != nil {
		panic(err)
	}

	spew.Dump(resp)

	fReq := &proto.ReqFetchEvents{
		FromDate: "16-3-2021",
	}

	fResp, err := client.FetchEvents(context.Background(), fReq)
	if err != nil {
		panic(err)
	}

	spew.Dump(fResp)

}
