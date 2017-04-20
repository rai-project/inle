package execution

type Error string

const (
	ErrContinue Error = "<continue input>"
	ErrQuit     Error = "<quit session>"
)

func (e Error) Error() string {
	return string(e)
}
