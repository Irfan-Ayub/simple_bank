package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("length must be between %d to %d", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	err := ValidateString(value, 3, 50)
	if err != nil {
		return fmt.Errorf("invalid username: %w", err)
	}

	if !isValidUsername(value) {
		return fmt.Errorf("username can only contain lowercase letters, numbers and underscore")
	}

	return nil
}

func ValidateFullname(value string) error {
	err := ValidateString(value, 3, 100)
	if err != nil {
		return fmt.Errorf("invalid username: %w", err)
	}

	if !isValidFullname(value) {
		return fmt.Errorf("fullname can only contain letters and spaces")
	}

	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 254); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}
	return nil
}
