package errorcodes

import "net/http"

type ErrorCode uint

const (
	NoError             ErrorCode = 0
	InternalServerError ErrorCode = 01
)

func (e ErrorCode) Message() string {
	messages := map[ErrorCode]string{
		NoError:             "",
		InternalServerError: "internal server error",
	}
	m := messages[e]
	return m
}

func (e ErrorCode) HttpStatusCode() int {
	httpErrorStatus := map[ErrorCode]int{
		NoError:             http.StatusOK,
		InternalServerError: http.StatusInternalServerError,
	}

	hes := httpErrorStatus[e]
	return hes
}
