package rpcserver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	calendarrpcapi "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/rpc/rpcapi"
	utils "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/rpc/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RPCServer struct {
	calendarrpcapi.UnimplementedApplicationServer
	logger interfaces.Logger
	app    interfaces.Application
	mu     sync.Mutex
	server *grpc.Server
}

func NewRPCServer(app interfaces.Application, logger interfaces.Logger) *RPCServer {
	return &RPCServer{
		app:    app,
		logger: logger,
	}
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
	if ident == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"идентификатор события не указан",
		)
	}
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
	events, err := s.app.ListEvents(stream.Context())
	if err != nil {
		return err
	}
	for i := range events {
		pbEvent := utils.EventToGRpcEvent(&events[i])

		if err := stream.Send(pbEvent); err != nil {
			return err
		}
	}
	return nil
}

func (s *RPCServer) ListNotSheduledEvents(_ *emptypb.Empty,
	stream calendarrpcapi.Application_ListNotSheduledEventsServer,
) error {
	events, err := s.app.ListNotSheduledEvents(stream.Context())
	if err != nil {
		return err
	}

	for i := range events {
		pbEvent := utils.EventToGRpcEvent(&events[i])

		if err := stream.Send(pbEvent); err != nil {
			return err
		}
	}

	return nil
}

func LoggingdUnaryInterceptor(
	logger interfaces.Logger,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		logger.Debug(
			"Получено gRPC сообщение: method=%q request={%v}",
			info.FullMethod,
			req,
		)

		start := time.Now()
		res, err := handler(ctx, req)
		duration := time.Since(start)
		code := status.Code(err)

		if err != nil {
			logger.Info("Выполнение метода: method=%q out={%v} code=%s duration=%s error=%v",
				info.FullMethod, res, code, duration, err)

			return res, utils.AppErrToGRPCError(err)
		}

		logger.Info("Выполнение метода: method=%q out={%v} code=%s duration=%s",
			info.FullMethod, res, code, duration)
		return res, nil
	}
}

type loggingServerStream struct {
	grpc.ServerStream

	logger interfaces.Logger
	method string

	received atomic.Uint64
	sent     atomic.Uint64
}

func (s *loggingServerStream) RecvMsg(message any) error {
	err := s.ServerStream.RecvMsg(message)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			s.logger.Error(
				"Ошибка чтения gRPC сообщения: method=%q error=%v",
				s.method,
				err,
			)
		}

		return err
	}

	s.received.Add(1)

	s.logger.Debug(
		"Получено gRPC сообщение: method=%q request={%v}",
		s.method,
		message,
	)

	return nil
}

func (s *loggingServerStream) SendMsg(message any) error {
	if err := s.ServerStream.SendMsg(message); err != nil {
		s.logger.Error(
			"Ошибка отправки gRPC сообщения: method=%q error=%v",
			s.method,
			err,
		)

		return err
	}

	s.sent.Add(1)

	s.logger.Debug(
		"Отправлено gRPC сообщение: method=%q response={%v}",
		s.method,
		message,
	)

	return nil
}

func LoggingStreamInterceptor(
	logger interfaces.Logger,
) grpc.StreamServerInterceptor {
	return func(
		srv any,
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		logger.Info(
			"Начало gRPC потока: method=%q client_stream=%t server_stream=%t",
			info.FullMethod,
			info.IsClientStream,
			info.IsServerStream,
		)

		wrappedStream := &loggingServerStream{
			ServerStream: stream,
			logger:       logger,
			method:       info.FullMethod,
		}

		start := time.Now()

		err := handler(srv, wrappedStream)

		duration := time.Since(start)
		code := status.Code(err)

		if err != nil {
			logger.Info(
				"Конец gRPC потока: method=%q code=%s duration=%s error=%v",
				info.FullMethod,
				code,
				duration,
				err,
			)

			return utils.AppErrToGRPCError(err)
		}

		logger.Info(
			"Конец gRPC потока: method=%q code=%s duration=%s",
			info.FullMethod,
			code,
			duration,
		)

		return nil
	}
}

func (s *RPCServer) getServer() *grpc.Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.server
}

func (s *RPCServer) Start(ctx context.Context, address string) error {
	listenConfig := net.ListenConfig{}
	listener, err := listenConfig.Listen(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf(
			"не удалось открыть gRPC listener %q: %w",
			address,
			err,
		)
	}
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(LoggingdUnaryInterceptor(s.logger)),
		grpc.StreamInterceptor(LoggingStreamInterceptor(s.logger)),
	)
	calendarrpcapi.RegisterApplicationServer(gRPCServer, s)
	s.mu.Lock()

	if s.server != nil {
		s.mu.Unlock()
		_ = listener.Close()

		return errors.New("gRPC сервер уже запущен")
	}

	s.server = gRPCServer
	s.mu.Unlock()

	// Очищаем поле, когда Serve завершится.
	defer func() {
		s.mu.Lock()

		// Проверка нужна, чтобы случайно не очистить ссылку
		// на другой экземпляр сервера.
		if s.server == gRPCServer {
			s.server = nil
		}

		s.mu.Unlock()
	}()

	s.logger.Info("Запуск gRPC сервера: address=%q", address)

	err = gRPCServer.Serve(listener)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("ошибка работы gRPC сервера: %w", err)
	}

	return nil
}

func (s *RPCServer) Stop() {
	server := s.getServer()

	if server == nil {
		s.logger.Info("gRPC сервер не запущен")
		return
	}

	s.logger.Info("Принудительный останов gRPC сервера")
	server.Stop()
}

func (s *RPCServer) GracefulStop() {
	server := s.getServer()

	if server == nil {
		s.logger.Info("gRPC сервер не запущен")
		return
	}

	s.logger.Info("Штатный останов gRPC сервера")
	server.GracefulStop()
}
