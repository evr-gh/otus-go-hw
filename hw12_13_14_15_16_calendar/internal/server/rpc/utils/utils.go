package utils

import (
	models "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/models"
	rpcapi "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/rpc/rpcapi"

	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func EventToGRpcEvent(event *models.Event) *rpcapi.Event {
	pbEvent := new(rpcapi.Event)
	pbEvent.Id = int32(event.ID)
	pbEvent.Title = event.Title
	pbEvent.Description = event.Description
	pbEvent.Time = timestamppb.New(event.Time)
	pbEvent.Duration = durationpb.New(event.Duration)
	pbEvent.Owner = event.Owner
	pbEvent.Notifyleadtime = durationpb.New(event.NotifyLeadTime)
	pbEvent.Sheduled = event.Sheduled
	return pbEvent
}

func GRpcEventToEvent(pbEvent *rpcapi.Event) *models.Event {
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
