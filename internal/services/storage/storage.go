package storage

import "errors"

var (
	ErrChannelExitsts     = errors.New("channel already exists")
	ErrChannelNotFound    = errors.New("channel not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrRowsIteration      = errors.New("rows iteration failed")
	ErrScanFailed         = errors.New("scan failed")
	ErrQueryFailed        = errors.New("query failed")
)
