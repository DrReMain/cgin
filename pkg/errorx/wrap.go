package errorx

import (
	"github.com/pkg/errors"

	"github.com/DrReMain/cgin/pkg/encoding/json"
)

var Is = errors.Is

func As(err error) (*ErrorX, bool) {
	if err != nil {
		var nerr *ErrorX
		if errors.As(err, &nerr) {
			return nerr, true
		}
	}
	return nil, false
}

// Parse convert string to ErrorX
// if string is a ErrorX marshal string, transform it
// else create a new ErrorX, and Detail is this string
func Parse(err string) *ErrorX {
	nerr := new(ErrorX)
	e := json.Unmarshal([]byte(err), nerr)
	if e != nil {
		nerr.Detail = err
	}
	return nerr
}

func FromError(err error) *ErrorX {
	if err == nil {
		return nil
	}
	if nerr, ok := As(err); ok {
		return nerr
	}

	return Parse(err.Error())
}
