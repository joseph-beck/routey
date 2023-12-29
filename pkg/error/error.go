package error

import "errors"

type ErrorCode int

const (
	DefaultErrorCode  ErrorCode = iota // 0
	NilErrorCode                       // 1
	RenderErrorCode                    // 2
	ServerErrorCode                    // 3
	QueryErrorCode                     // 4
	RedirectErrorCode                  // 5
	HTMLErrorCode                      // 6
	NoDataErrorCode                    // 7
)

type Error struct {
	Message string
	Error   error
	Code    ErrorCode
}

var (
	NilError = Error{
		Message: "",
		Error:   nil,
		Code:    NilErrorCode,
	}
	RenderError = Error{
		Message: "Render Error Occurred",
		Error:   errors.New("render error occurred"),
		Code:    RenderErrorCode,
	}
	ServerError = Error{
		Message: "Server Error Occurred",
		Error:   errors.New("server error occurred"),
		Code:    ServerErrorCode,
	}
	QueryError = Error{
		Message: "Query Error Occurred",
		Error:   errors.New("query error occurred"),
		Code:    QueryErrorCode,
	}
	RedirectError = Error{
		Message: "Redirect Error Occurred",
		Error:   errors.New("redirect error occurred"),
		Code:    RedirectErrorCode,
	}
	HTMLError = Error{
		Message: "HTML Error Occurred",
		Error:   errors.New("html error occurred"),
		Code:    HTMLErrorCode,
	}
	NoDataError = Error{
		Message: "No Data Error Occurred",
		Error:   errors.New("no data error occurred"),
		Code:    NoDataErrorCode,
	}
)

// Creates a new Error
func New(m string, e error) Error {
	return Error{
		Message: m,
		Error:   e,
		Code:    DefaultErrorCode,
	}
}

// Get the message of the Error
func (e Error) String() string {
	return e.Message
}

// Is this Error of the same code as another?
func (e Error) Equal(o Error) bool {
	return e.Code == o.Code
}
