package server

import (
	"context"
	"net"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	calendarrpcapi "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/rpc/rpcapi"
	utils "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/rpc/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RPCServer struct {
	calendarrpcapi.UnimplementedApplicationServer
	logger interfaces.Logger
	app    interfaces.Application
	server *grpc.Server
}

func NewRPCServer(app interfaces.Application, logger interfaces.Logger) *RPCServer {
	self := &RPCServer{}
	self.app = app
	self.logger = logger
	return self
}

func (s *RPCServer) CreateEvent(ctx context.Context, pbEvent *calendarrpcapi.Event) (*calendarrpcapi.Event, error) {
	event := utils.GRpcEventToEvent(pbEvent)
	createdEvent, err := s.app.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return utils.EventToGRpcEvent(createdEvent), nil
}

func (s *RPCServer) ReadEvent(ctx context.Context, ident *calendarrpcapi.Id) (*calendarrpcapi.Event, error) {
	event, err := s.app.ReadEvent(ctx, int(ident.Id))
	if err != nil {
		return nil, err
	}
	return utils.EventToGRpcEvent(event), nil
}

func (s *RPCServer) UpdateEvent(ctx context.Context, pbEvent *calendarrpcapi.Event) (*calendarrpcapi.Event, error) {
	event := utils.GRpcEventToEvent(pbEvent)
	updatedEvent, err := s.app.UpdateEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return utils.EventToGRpcEvent(updatedEvent), nil
}

func (s *RPCServer) DeleteEvent(ctx context.Context, pbEvent *calendarrpcapi.Event) (*calendarrpcapi.Event, error) {
	event := utils.GRpcEventToEvent(pbEvent)
	deletedEvent, err := s.app.DeleteEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return utils.EventToGRpcEvent(deletedEvent), nil
}

func (s *RPCServer) ListEvents(_ *emptypb.Empty, stream calendarrpcapi.Application_ListEventsServer) error {
	ctx := stream.Context()

	events, err := s.app.ListEvents(ctx)
	if err != nil {
		return err
	}
	for _, event := range events {
		tmp := event
		pbEvent := utils.EventToGRpcEvent(&tmp)
		if err := stream.Send(pbEvent); err != nil {
			return err
		}
	}
	return nil
}

func (s *RPCServer) ListNotSheduledEvents(_ *emptypb.Empty,
	stream calendarrpcapi.Application_ListNotSheduledEventsServer,
) error {
	ctx := stream.Context()
	events, err := s.app.ListNotSheduledEvents(ctx)
	if err != nil {
		return err
	}
	for _, event := range events {
		tmp := event
		pbEvent := utils.EventToGRpcEvent(&tmp)
		if err := stream.Send(pbEvent); err != nil {
			return err
		}
	}
	return nil
}

func LoggedUnaryInterceptor(logger interfaces.Logger) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Info("%q <-- OBJECT{%s}", info.FullMethod, req)
		return handler(ctx, req)
	}
}

func LoggedStreamInterceptor(logger interfaces.Logger) func(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		logger.Info("%q %s", info.FullMethod, info)
		return handler(srv, stream)
	}
}

func (s *RPCServer) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(LoggedUnaryInterceptor(s.logger)),
		grpc.StreamInterceptor(LoggedStreamInterceptor(s.logger)),
	)
	calendarrpcapi.RegisterApplicationServer(gRPCServer, s)
	s.server = gRPCServer
	s.logger.Info("Запуск gRPC сервера")
	return s.server.Serve(lis)
}

func (s *RPCServer) Stop() {
	s.logger.Info("Останов gRPC сервера")
	s.server.Stop()
}

func (s *RPCServer) GracefulStop() {
	s.logger.Info("Штатный останов gRPC сервера")
	s.server.GracefulStop()
}
