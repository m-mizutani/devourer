package types

type Error struct {
	msg string
}

func (x *Error) Error() string { return x.msg }

func (x *Error) Is(err error) bool {
	if e, ok := err.(*Error); ok {
		return x.msg == e.msg
	}
	return false
}

func newError(msg string) *Error {
	return &Error{
		msg: msg,
	}
}

var (
	ErrInvalidOption = newError("invalid option")
)
