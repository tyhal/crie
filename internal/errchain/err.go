// Package errchain keeps consistency of displaying error messages with more context.
package errchain

import "fmt"

// ErrorChain wraps an error and provides helpers to add contextual messages
// while preserving the original error using Go's error wrapping.
type ErrorChain struct {
	err error
}

// From creates a new ErrorChain rooted at the provided error so that
// additional context can be added in a consistent manner.
func From(err error) ErrorChain {
	return ErrorChain{err: err}
}

// ErrorF formats a contextual message with the given format and arguments,
// and returns the original error wrapped with that context.
func (ec ErrorChain) ErrorF(format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	return fmt.Errorf("%s ⟶ %w", msg, ec.err)
}

// Error adds a contextual message and returns the original error wrapped
// with that context.
func (ec ErrorChain) Error(msg string) error {
	return fmt.Errorf("%s ⟶ %w", msg, ec.err)
}
