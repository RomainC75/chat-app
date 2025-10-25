package custom_errors

import "fmt"

type ErrNotFound struct {
	err error
}

func NewErrNotFound(err error) ErrNotFound {
	return ErrNotFound{
		err: err,
	}
}

func (nfe ErrNotFound) Error() string {
	return fmt.Sprintf("Not Found Error: %w", nfe.err)
}
