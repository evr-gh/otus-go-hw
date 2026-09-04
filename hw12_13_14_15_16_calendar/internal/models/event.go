package data

import (
	"time"
)

type Event struct {
	ID             int
	Title          string        `json:"title"`
	Time           time.Time     `json:"startat"`
	Duration       time.Duration `json:"duration"`
	Description    string        `json:"description"`
	Owner          string        `json:"owner"`
	NotifyLeadTime time.Duration `json:"notifyleadtime"`
	Sheduled       bool          `json:"sheduled"`
}
