package cache

import (
	"errors"
)

var (
	// ErrNotFound record not found error
	ErrNotFound = errors.New("not found")
	// ErrDataCorruption value in redis not JSON or some failure in unmarshal JSON
	ErrDataCorruption = errors.New("data corruption")
)
