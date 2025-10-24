package linter

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFailedResultError_Error_ReturnsUnderlying(t *testing.T) {
	under := errors.New("linter reported failure")
	fre := &FailedResultError{err: under}

	// Ensure Error() returns the underlying message (and does not recurse)
	assert.Equal(t, under.Error(), fre.Error())
}

func TestResult_NilPassThrough(t *testing.T) {
	assert.Nil(t, Result(nil))
}

func TestResult_WrapsNonNilError(t *testing.T) {
	under := errors.New("failed rules")
	err := Result(under)

	// Not nil and of the correct type
	if assert.NotNil(t, err) {
		var fre *FailedResultError
		if assert.ErrorAs(t, err, &fre) {
			// Error string should match the underlying error
			assert.Equal(t, under.Error(), fre.Error())
		}
	}
}
