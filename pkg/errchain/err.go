// Package errchain keeps consistency of displaying error messages with more context.
package errchain

import "fmt"

// Chain represents a helper that wraps an existing error and allows
// adding contextual messages while preserving the original error using
// Go's error wrapping.
type Chain interface {
	LinkF(format string, a ...any) error
	Link(msg string) error
}

// errorChain wraps an error and provides helpers to add contextual messages
// while preserving the original error using Go's error wrapping.
type errorChain struct {
	err error
}

// From creates a new errorChain rooted at the provided error so that
// additional context can be added in a consistent manner.
func From(err error) Chain {
	return errorChain{err: err}
}

// LinkF formats a contextual message with the given format and arguments,
// and returns the original error wrapped with that context.
func (ec errorChain) LinkF(format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	return fmt.Errorf("%s ⟶ %w", msg, ec.err)
}

// Link adds a contextual message and returns the original error wrapped
// with that context.
func (ec errorChain) Link(msg string) error {
	return fmt.Errorf("%s ⟶ %w", msg, ec.err)
}
