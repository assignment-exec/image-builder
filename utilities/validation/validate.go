// Package validation contains utilities to help run validators.
package validation

import "github.com/pkg/errors"

// Validator is a function that performs any kind of validation.
type Validator func() error

// Validate validates one or more validators.
// It returns any error encountered wrapped around
// with the base error message.
func Validate(errorMsg string, validators ...Validator) error {
	for _, v := range validators {
		if err := v(); err != nil {
			return errors.Wrap(err, errorMsg)
		}
	}
	return nil
}
