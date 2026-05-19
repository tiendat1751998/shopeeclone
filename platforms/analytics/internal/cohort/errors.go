package cohort

import "errors"

var (
	ErrCohortInvalid = errors.New("cohort: invalid cohort definition")
	ErrCohortNotFound = errors.New("cohort: cohort not found")
)
