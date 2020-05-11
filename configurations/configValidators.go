package configurations

import (
	"assignment-exec/image-builder/utilities/validation"
	"github.com/pkg/errors"
)

type configValidator func(AssignmentEnvConfig) error

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

func withBaseImageValidator() configValidator {
	return func(cfg AssignmentEnvConfig) error {
		// Base Image name cannot be empty string.
		if cfg.BaseImage == "" {
			return errors.New("base image name cannot be empty string")
		}

		return validateBaseImage(cfg.BaseImage)
	}
}

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
