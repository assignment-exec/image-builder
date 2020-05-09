package configurations

import (
	"assignment-exec/image-builder/utilities/validation"
	"github.com/pkg/errors"
)

type dockerConfigValidator func(DockerConfig) error

func ValidatorForConfig(d DockerConfig, configValidators ...dockerConfigValidator) validation.Validator {
	return func() error {
		for _, cv := range configValidators {
			if err := cv(d); err != nil {
				return err
			}
		}
		return nil
	}
}

func withBaseImageValidator() dockerConfigValidator {
	return func(d DockerConfig) error {
		// Base Image name cannot be empty string.
		if d.BaseImage == "" {
			return errors.New("base image name cannot be empty string")
		}

		err := validateBaseImage(d.BaseImage)
		if err != nil {
			return err
		}
		return nil
	}
}

func withLanguageValidator() dockerConfigValidator {
	return func(d DockerConfig) error {
		// Language name and version name cannot be empty string.
		if d.Dependencies.Language.Name == "" || d.Dependencies.Language.Version == "" {
			return errors.New("language name and version cannot be empty string")
		}

		lang := d.Dependencies.Language.Name
		version := d.Dependencies.Language.Version
		if err := validateLang(lang, version); err != nil {
			return errors.Wrap(err, "programming language not supported")
		}
		return nil
	}
}

func withLibsValidator() dockerConfigValidator {
	return func(d DockerConfig) error {
		// Library installation commands cannot be empty strings.
		for s, libInstallCmd := range d.Dependencies.Libraries {
			if s == "" || libInstallCmd.Cmd == "" {
				return errors.New("library installation command cannot be empty string")
			}
		}
		return nil
	}
}
