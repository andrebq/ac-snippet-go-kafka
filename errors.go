package main

import "net/http"

type (
	statusErr struct {
		msg    string
		status int
	}
)

func (se statusErr) Error() string {
	return se.msg
}

func (se statusErr) Status() int {
	return se.status
}

func (se statusErr) HTTPStatus() int {
	return se.Status()
}

func badRequest(msg string) error {
	return statusErr{
		msg:    msg,
		status: http.StatusBadRequest,
	}
}

func badGateway(msg string) error {
	return statusErr{
		msg:    msg,
		status: http.StatusBadGateway,
	}
}
