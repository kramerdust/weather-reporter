package forecaster

import "net/http"

func IsNotFound(err error) bool {
	if fcErr, ok := err.(*fcError); ok {
		return fcErr.code == http.StatusNotFound
	} else {
		return false
	}
}

type fcError struct {
	msg  string
	code int
}

func (err *fcError) Error() string {
	return err.msg
}
