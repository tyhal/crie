//go:build !linux
// +build !linux

package cli

// MaxConcurrency returns the max concurrency name
func (e *Lint) MaxConcurrency() int {
	return 128
}
