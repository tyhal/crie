package linter

// FailedResultError means a linter didn't error but returned a failed result
type FailedResultError struct {
	err error
}

// Error implements the error interface
func (e *FailedResultError) Error() string {
	return e.err.Error()
}

// Result conditionally wraps an error with a FailedResultError or otherwise passes through nil, it should be used when a linter didn't error but returned a failed result
func Result(err error) error {
	if err == nil {
		return nil
	}
	return &FailedResultError{err: err}
}
