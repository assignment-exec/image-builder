// Package configurations provides primitives to read and store the
// assignment environment configuration yaml file, get the docker instructions
// in the specific format for every configuration.
package configurations

import (
	"assignment-exec/image-builder/utilities/validation"
	"github.com/pkg/errors"
)

// configValidator is a function interface that is supplied
// as different options while validating parameters of `AssignmentEnvConfig`.
type configValidator func(AssignmentEnvConfig) error

// ValidatorForConfig takes the `AssignmentEnvConfig` to be validated
// and one or more validator options to validate all parameters.
func ValidatorForConfig(cfg AssignmentEnvConfig, configValidators ...configValidator) validation.Validator {
	return func() error {
		for _, cfgValidator := range configValidators {
			if err := cfgValidator(cfg); err != nil {
				return err
			}
		}
		return nil
	}
}

// withBaseImageValidator is used as an option while validating. It returns
// `configValidator` function. This returned function in turn validates the base
// image provided in assignment environment config.
func withBaseImageValidator() configValidator {
	return func(cfg AssignmentEnvConfig) error {
		// Base Image name cannot be empty string.
		if cfg.BaseImage == "" {
			return errors.New("base image name cannot be empty string")
		}

		return validateBaseImage(cfg.BaseImage)
	}
}

// withLanguageValidator is used as an option while validating. It returns
// `configValidator` function. This returned function in turn validates the
// language information provided in assignment environment config.
func withLanguageValidator() configValidator {
	return func(cfg AssignmentEnvConfig) error {
		// Language name and version name cannot be empty string.
		if cfg.Deps.Language.Name == "" || cfg.Deps.Language.Version == "" {
			return errors.New("language name and version cannot be empty string")
		}

		lang := cfg.Deps.Language.Name
		version := cfg.Deps.Language.Version
		if err := validateLang(lang, version); err != nil {
			return errors.Wrap(err, "programming language not supported")
		}
		return nil
	}
}

// withLibsValidator is used as an option while validating. It returns
// `configValidator` function. This returned function in turn validates the
// library names and their respective installation commands provided in
// assignment environment config.
func withLibsValidator() configValidator {
	return func(cfg AssignmentEnvConfig) error {
		// Library installation commands cannot be empty strings.
		for s, libInstallCmd := range cfg.Deps.Libraries {
			if s == "" || libInstallCmd.Cmd == "" {
				return errors.New("library installation command cannot be empty string")
			}
		}
		return nil
	}
}
