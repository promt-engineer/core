package history

import "errors"

var (
	ErrSpinNotFound        = errors.New("spin not found")
	ErrIsDemoRequiredField = errors.New("is demo required field")
)
