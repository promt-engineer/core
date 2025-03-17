package validator

import (
	"fmt"
	"strings"
)

type sliceValidateError []error

func (err sliceValidateError) Error() string {
	errMsgs := []string{}

	for i, e := range err {
		if e == nil {
			continue
		}

		errMsgs = append(errMsgs, fmt.Sprintf("[%d]: %s", i, e.Error()))
	}

	return strings.Join(errMsgs, "\n")
}
