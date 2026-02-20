package gofp

type ResultError string
type resultPanic struct{ err error }

func (e ResultError) Error() string {
	return string(e)
}

const ErrNoResults ResultError = "result: no results provided"
