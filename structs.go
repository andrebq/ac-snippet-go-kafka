package main

import (
	"time"
)

type (
	Message struct {
		ID          string
		Destination string
		Sender      string
		Text        string
		LocalTime   time.Time
		ServerTime  time.Time
	}

	key byte
)

var (
	authKey = key(0)
)
