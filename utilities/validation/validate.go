package validation

import "github.com/pkg/errors"

type Validator func() error

func Validate(baseErrMsg string, validators ...Validator) error {
	for _, v := range validators {
		if err := v(); err != nil {
			return errors.Wrap(err, baseErrMsg)
		}
	}

	return nil
}
