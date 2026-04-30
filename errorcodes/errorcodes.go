package errorcodes

import "net/http"

type ErrorCode uint

const (
	NoError             ErrorCode = 0
	InternalServerError ErrorCode = 1
	BadRequest          ErrorCode = 2
	Unauthorized        ErrorCode = 3
	NotFound            ErrorCode = 4
	Forbidden           ErrorCode = 5
	Conflict            ErrorCode = 6
)

func (e ErrorCode) Message() string {
	messages := map[ErrorCode]string{
		NoError:             "",
		InternalServerError: "internal server error",
		BadRequest:          "bad request",
		Unauthorized:        "unauthorized access",
		NotFound:            "resource not found",
		Forbidden:           "forbidden",
		Conflict:            "conflict",
	}
	return messages[e]
}

func (e ErrorCode) HttpStatusCode() int {
	httpErrorStatus := map[ErrorCode]int{
		NoError:             http.StatusOK,
		InternalServerError: http.StatusInternalServerError,
		BadRequest:          http.StatusBadRequest,
		Unauthorized:        http.StatusUnauthorized,
		NotFound:            http.StatusNotFound,
		Forbidden:           http.StatusForbidden,
		Conflict:            http.StatusConflict,
	}
	return httpErrorStatus[e]
}
