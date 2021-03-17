package api

import (
	"github.com/ark1790/ch/eventstore/belt"
	"github.com/ark1790/ch/eventstore/proto"
	"github.com/ark1790/ch/eventstore/repo"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	proto.UnimplementedEventStoreServer
	eventRepo repo.EventRepo
	queue     belt.Queue
}

var logger = logrus.NewEntry(logrus.New())

var unaryInterceptors = []grpc.UnaryServerInterceptor{
	grpc_logrus.UnaryServerInterceptor(logger),
}

func NewServer(er repo.EventRepo, q belt.Queue) *Server {
	srvr := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(unaryInterceptors...),
		),
	)

	s := &Server{
		Server:    srvr,
		eventRepo: er,
		queue:     q,
	}

	proto.RegisterEventStoreServer(srvr, s)

	return s
}
