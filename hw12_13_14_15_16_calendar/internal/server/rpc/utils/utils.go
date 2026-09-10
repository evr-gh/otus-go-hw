package utils

import (
	"errors"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	models "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/models"
	rpcapi "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/rpc/rpcapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func EventToGRpcEvent(event *models.Event) *rpcapi.Event {
	if event != nil {
		pbEvent := new(rpcapi.Event)
		pbEvent.Id = int64(event.ID)
		pbEvent.Title = event.Title
		pbEvent.Description = event.Description
		pbEvent.Time = timestamppb.New(event.Time)
		pbEvent.Duration = durationpb.New(event.Duration)
		pbEvent.Owner = event.Owner
		pbEvent.Notifyleadtime = durationpb.New(event.NotifyLeadTime)
		pbEvent.Sheduled = event.Sheduled
		return pbEvent
	}
	return nil
}

func GRpcEventToEvent(pbEvent *rpcapi.Event) *models.Event {
	if pbEvent != nil {
		event := new(models.Event)
		event.ID = int(pbEvent.Id)
		event.Title = pbEvent.Title
		event.Description = pbEvent.Description
		event.Time = pbEvent.Time.AsTime()
		event.Duration = pbEvent.Duration.AsDuration()
		event.Owner = pbEvent.Owner
		event.NotifyLeadTime = pbEvent.Notifyleadtime.AsDuration()
		event.Sheduled = pbEvent.Sheduled
		return event
	}
	return nil
}

func AppErrToGRPCError(err error) error {
	switch {
	case err == nil:
		return nil

	case errors.Is(err, interfaces.ErrNoData):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, interfaces.ErrNoEvent):
		return status.Error(codes.InvalidArgument, err.Error())

	default:
		return status.Error(codes.Internal, "внутренняя ошибка сервера")
	}
}
