package errchain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_WrapsAndAddsContext(t *testing.T) {
	base := errors.New("base failure")

	err := From(base).Error("doing important work")

	// It should wrap the base error
	assert.ErrorIs(t, err, base)

	// The message should contain our context and the base message
	s := err.Error()
	assert.Contains(t, s, "doing important work")
	assert.Contains(t, s, base.Error())
}

func TestErrorF_WrapsAndFormats(t *testing.T) {
	base := errors.New("not found")

	err := From(base).ErrorF("fetching id %d", 42)

	assert.ErrorIs(t, err, base)

	s := err.Error()
	assert.Contains(t, s, "fetching id 42")
	assert.Contains(t, s, base.Error())
}

func TestChaining_PreservesOriginalAndAllContexts(t *testing.T) {
	base := errors.New("disk full")

	// First wrap
	err1 := From(base).Error("writing cache")
	// Second wrap by starting a new chain from the previous error
	err2 := From(err1).ErrorF("attempt %d", 3)

	assert.ErrorIs(t, err2, base)

	s := err2.Error()
	assert.Contains(t, s, "writing cache")
	assert.Contains(t, s, "attempt 3")
}
