package errors

import (
	"errors"
)

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrRecordNotUpdated = errors.New("record not updated")
	ErrRecordNotCreated = errors.New("record not created")
	ErrRecordNotDeleted = errors.New("record not deleted")
)
