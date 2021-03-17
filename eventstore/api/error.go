package api

import (
	"net/http"

	"github.com/ark1790/ch/eventstore/proto"
	protobuf "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type validationError map[string][]string

func (v validationError) Error() *errdetails.BadRequest {
	details := []*errdetails.BadRequest_FieldViolation{}
	for k, v := range v {
		if len(v) > 0 {
			detail := &errdetails.BadRequest_FieldViolation{
				Field:       k,
				Description: v[0],
			}

			details = append(details, detail)
		}
	}

	return &errdetails.BadRequest{
		FieldViolations: details,
	}
}

func (v validationError) add(key string, val string) {
	v[key] = append(v[key], val)
}

func (v validationError) extend(prefix string, err *validationError) {
	if err == nil {
		return
	}
	for k, e := range *err {
		k = prefix + k
		v[k] = append(v[k], e...)
	}
}

type errorCode struct {
	Code   string
	Status int
}

var (
	errInvalidData         = &errorCode{Code: "422001", Status: http.StatusUnprocessableEntity}
	errNotFound            = &errorCode{Code: "404001", Status: http.StatusNotFound}
	errInternalServerError = &errorCode{Code: "500001", Status: http.StatusInternalServerError}
)

func newAPIError(title string, ec *errorCode, a []*any.Any) *proto.Errors {
	return &proto.Errors{
		Id:      uuid.NewV4().String(),
		Status:  int32(ec.Status),
		Code:    ec.Code,
		Title:   title,
		Details: a,
	}
}

func marshalAny(details ...protobuf.Message) ([]*any.Any, error) {
	any := []*any.Any{}
	for _, detail := range details {
		a, err := ptypes.MarshalAny(detail)
		if err != nil {
			return nil, err
		}
		any = append(any, a)
	}

	return any, nil
}
