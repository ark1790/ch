package api

import (
	"context"
	"strings"
	"time"

	"github.com/ark1790/ch/eventstore/model"
	"github.com/ark1790/ch/eventstore/proto"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func validateCreateEvent(pld *proto.ReqCreateEvent) error {
	pld.Email = strings.TrimSpace(pld.Email)
	pld.Environment = strings.TrimSpace(pld.Environment)
	pld.Message = strings.TrimSpace(pld.Message)
	pld.Component = strings.TrimSpace(pld.Component)

	errV := validationError{}
	if pld.Email == "" {
		errV.add("email", "is required")
	} else if !govalidator.IsEmail(pld.Email) {
		errV.add("email", "is not valid")
	}
	if pld.Environment == "" {
		errV.add("environment", "is required")
	}
	if pld.Message == "" {
		errV.add("message", "is required")
	}
	if pld.Component == "" {
		errV.add("component", "is required")
	}

	if len(errV) > 0 {
		st := status.New(codes.InvalidArgument, "Validation Error")
		eAny, err := marshalAny(errV.Error())
		if err != nil {
			return st.Err()
		}

		ds, err := st.WithDetails(
			newAPIError("Invalid Data", errInvalidData, eAny),
		)
		if err != nil {
			return ds.Err()
		}
		return ds.Err()
	}

	return nil
}

func (s *Server) CreateEvent(ctx context.Context, r *proto.ReqCreateEvent) (*proto.RespCreateEvent, error) {
	if err := validateCreateEvent(r); err != nil {
		return nil, err
	}

	data := r.Data.AsMap()

	evt := &model.Event{
		Email:       r.Email,
		Environment: r.Environment,
		Component:   r.Component,
		Message:     r.Message,
		Data:        data,
		CreatedAt:   time.Unix(r.CreatedAt, 0),
	}

	evt.Pre()

	// publish to belt

	msg, err := evt.Marshal()
	if err != nil {
		return nil, err
	}

	s.queue.Push(msg)

	resp := &proto.RespCreateEvent{
		Event: toProtoEvent(*evt),
	}

	return resp, err
}

func validateFetchEvents(pld *proto.ReqFetchEvents) error {

	pld.Id = strings.TrimSpace(pld.Id)
	pld.Email = strings.TrimSpace(pld.Email)
	pld.Environment = strings.TrimSpace(pld.Environment)
	pld.Message = strings.TrimSpace(pld.Message)
	pld.Component = strings.TrimSpace(pld.Component)

	if pld.Page <= 0 {
		pld.Page = 1
	}
	if pld.PerPage <= 0 || pld.PerPage > 20 {
		pld.PerPage = 20
	}

	errV := validationError{}
	if pld.Email != "" && !govalidator.IsEmail(pld.Email) {
		errV.add("email", "is not valid")
	}
	if pld.FromDate != "" {
		_, err := time.Parse("2-1-2006", pld.FromDate)
		if err != nil {
			errV.add("from_date", "is not valid")
		}
	}

	if len(errV) > 0 {
		st := status.New(codes.InvalidArgument, "Validation Error")
		eAny, err := marshalAny(errV.Error())
		if err != nil {
			return st.Err()
		}

		ds, err := st.WithDetails(
			newAPIError("Invalid Data", errInvalidData, eAny),
		)
		if err != nil {
			return ds.Err()
		}
		return ds.Err()
	}

	return nil
}

func transformQuery(r *proto.ReqFetchEvents) map[string]string {
	qry := map[string]string{}
	if r.Id != "" {
		qry["id"] = r.Id
	}
	if r.Email != "" {
		qry["email"] = r.Email
	}
	if r.Environment != "" {
		qry["environment"] = r.Environment
	}
	if r.Component != "" {
		qry["component"] = r.Component
	}
	if r.Message != "" {
		qry["message"] = r.Message
	}

	if r.From != 0 {
		from := time.Unix(r.From, 0)
		qry["from"] = from.Format(time.RFC3339)
	}
	if r.To != 0 {
		to := time.Unix(r.To, 0)
		qry["to"] = to.Format(time.RFC3339)
	}
	if r.CreatedAt != 0 {
		createdAt := time.Unix(r.CreatedAt, 0)
		qry["created_at"] = createdAt.Format(time.RFC3339)
	}

	if r.FromDate != "" {
		from, _ := time.Parse("2-1-2006", r.FromDate)
		qry["from"] = from.Format(time.RFC3339)

	}

	return qry
}

func toProtoEvent(e model.Event) *proto.Event {

	pCAt := e.CreatedAt.Unix()
	msData, _ := e.Data.(map[string]interface{})

	d, _ := structpb.NewStruct(msData)

	return &proto.Event{
		Id:          e.ID,
		Email:       e.Email,
		Environment: e.Environment,
		Component:   e.Component,
		Message:     e.Message,
		Data:        d,
		CreatedAt:   pCAt,
	}
}

func (s *Server) FetchEvents(ctx context.Context, r *proto.ReqFetchEvents) (*proto.RespFetchEvents, error) {
	if err := validateFetchEvents(r); err != nil {

		return nil, err
	}

	offset := int((r.Page - 1) * r.PerPage)
	limit := int(r.PerPage)

	qry := transformQuery(r)

	mEvts, err := s.eventRepo.FetchEvents(ctx, qry, offset, limit)
	if err != nil {
		return nil, err
	}

	evts := []*proto.Event{}
	for _, e := range mEvts {
		evts = append(evts, toProtoEvent(e))
	}

	resp := &proto.RespFetchEvents{
		Events:  evts,
		Page:    r.Page,
		PerPage: r.PerPage,
	}

	return resp, nil
}
