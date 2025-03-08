package service

import "errors"

var (
	ErrInvalidParams   = errors.New("invalid params")
	ErrVectorTransform = errors.New("vector transform failed")
	ErrDataNotFound    = errors.New("data not found")
)
