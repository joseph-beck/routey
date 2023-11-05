package error

import "errors"

type Error struct {
	Message string
	Error   error
}

var (
	NilError = Error{
		Message: "",
		Error:   nil,
	}
	RenderError = Error{
		Message: "Render Error Occurred",
		Error:   errors.New("render error occurred"),
	}
	ServerError = Error{
		Message: "Server Error Occurred",
		Error:   errors.New("server error occurred"),
	}
)

func New(m string, e error) Error {
	return Error{
		Message: m,
		Error:   e,
	}
}

func (e Error) String() string {
	return e.Message
}
