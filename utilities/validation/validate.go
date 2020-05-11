package validation

import "github.com/pkg/errors"

type Validator func() error

func Validate(errorMsg string, validators ...Validator) error {
	for _, v := range validators {
		if err := v(); err != nil {
			return errors.Wrap(err, errorMsg)
		}
	}
	return nil
}
