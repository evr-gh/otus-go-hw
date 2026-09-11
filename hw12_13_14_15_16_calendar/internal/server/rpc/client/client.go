package client

import (
	"context"
	"errors"
	"io"
	"log"

	calendarrpcapi "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/rpc/rpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	grpcClient calendarrpcapi.ApplicationClient
	connection *grpc.ClientConn
}

var localOpts = []grpc.CallOption{}

func (c *Client) Connect(dsn string) error {
	connection, err := grpc.NewClient(dsn, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed on client: %v", err)
		return err
	}
	c.connection = connection
	c.grpcClient = calendarrpcapi.NewApplicationClient(c.connection)
	return nil
}

func (c *Client) Close() error {
	err := c.connection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateEvent(ctx context.Context, event *calendarrpcapi.Event) (*calendarrpcapi.Event, error) {
	return c.grpcClient.CreateEvent(ctx, event, localOpts...)
}

func (c *Client) ReadEvent(ctx context.Context, event *calendarrpcapi.Id) (*calendarrpcapi.Event, error) {
	return c.grpcClient.ReadEvent(ctx, event, localOpts...)
}

func (c *Client) UpdateEvent(ctx context.Context, event *calendarrpcapi.Event) (*calendarrpcapi.Event, error) {
	return c.grpcClient.UpdateEvent(ctx, event, localOpts...)
}

func (c *Client) DeleteEvent(ctx context.Context, event *calendarrpcapi.Event) (*calendarrpcapi.Event, error) {
	return c.grpcClient.DeleteEvent(ctx, event, localOpts...)
}

func (c *Client) ListEvents(ctx context.Context) ([]*calendarrpcapi.Event, error) {
	events, err := c.grpcClient.ListEvents(ctx, &emptypb.Empty{}, localOpts...)
	if err != nil {
		return nil, err
	}
	pbEvents := make([]*calendarrpcapi.Event, 0)
	for {
		pbEvent, err := events.Recv()
		// err := events.RecvMsg(&pbEvent)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		pbEvents = append(pbEvents, pbEvent)
	}
	return pbEvents, nil
}

func (c *Client) ListNotSheduledEvents(ctx context.Context) ([]*calendarrpcapi.Event, error) {
	events, err := c.grpcClient.ListNotSheduledEvents(ctx, &emptypb.Empty{}, localOpts...)
	if err != nil {
		return nil, err
	}
	pbEvents := make([]*calendarrpcapi.Event, 0)
	for {
		pbEvent, err := events.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		pbEvents = append(pbEvents, pbEvent)
	}
	return pbEvents, nil
}
