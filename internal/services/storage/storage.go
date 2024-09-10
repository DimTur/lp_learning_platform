package storage

import "errors"

var (
	ErrChannelExitsts  = errors.New("channel already exists")
	ErrChannelNotFound = errors.New("channel not found")
)
