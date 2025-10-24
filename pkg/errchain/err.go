// Package errchain keeps consistency of displaying error messages with more context.
package errchain

import "fmt"

type ErrorChain struct {
	err error
}

func From(err error) ErrorChain {
	return ErrorChain{err: err}
}

func (ec ErrorChain) ErrorF(format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	return fmt.Errorf("%s ⟶ %w", msg, ec.err)
}

func (ec ErrorChain) Error(msg string) error {
	return fmt.Errorf("%s ⟶ %w", msg, ec.err)
}
