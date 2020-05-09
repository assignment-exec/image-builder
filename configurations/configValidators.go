package configurations

import (
	"assignment-exec/image-builder/utilities/validation"
	"github.com/pkg/errors"
)

type configValidator func(Config) error

func ValidatorForConfig(t Config, configValidators ...configValidator) validation.Validator {
	return func() error {
		for _, tv := range configValidators {
			if err := tv(t); err != nil {
				return err
			}
		}

		return nil
	}
}

func withBaseImageValidator() configValidator {
	return func(c Config) error {
		// Base Image name cannot be empty string.
		if c.BaseImage == "" {
			return errors.New("base image name cannot be empty string")
		}

		err := validateBaseImage(c.BaseImage)
		if err != nil {
			return err
		}
		return nil
	}
}

func withLanguageValidator() configValidator {
	return func(c Config) error {
		// Language name and version name cannot be empty string.
		if c.Deps.Language.Name == "" || c.Deps.Language.Version == "" {
			return errors.New("language name and version cannot be empty string")
		}

		lang := c.Deps.Language.Name
		version := c.Deps.Language.Version
		if err := validateLang(lang, version); err != nil {
			return errors.Wrap(err, "programming language not supported")
		}
		return nil
	}
}

func withLibsValidator() configValidator {
	return func(c Config) error {
		// Library installation commands cannot be empty strings.
		for s, libInstallCmd := range c.Deps.Libs {
			if s == "" || libInstallCmd.Cmd == "" {
				return errors.New("library installation command cannot be empty string")
			}
		}
		return nil
	}
}
