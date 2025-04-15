package util

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

func ValidateEmail(email string) string {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return "Email address is not valid"
	}
	// TODO Add duplicate email check
	return ""
}

var usernameRegex = regexp.MustCompile(`^[a-zA-ZæøåÆØÅ0-9_.-]+$`)

func ValidateUsername(username string) string {
	if len(username) < 3 || len(username) > 30 {
		return "Username must be between 3 and 30 characters"
	}

	if !usernameRegex.MatchString(username) {
		return "Username can only contain Danish letters, digits, '.', '_', and '-'"
	}

	// TODO: Add duplicate username check
	return ""
}

func ValidatePassword(password string) string {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
		specials   = "!@#$%^&*()-_=+[]{}|;:,.<>/?`~"
	)

	if len(password) < 8 {
		return "Password must be at least 8 characters long"
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case strings.ContainsRune(specials, char):
			hasSpecial = true
		}
	}

	var requirements []string
	if !hasUpper {
		requirements = append(requirements, "one uppercase letter")
	}
	if !hasLower {
		requirements = append(requirements, "one lowercase letter")
	}
	if !hasNumber {
		requirements = append(requirements, "one number")
	}
	if !hasSpecial {
		requirements = append(requirements, "one special character")
	}

	if len(requirements) > 0 {
		return "Password must contain at least " + strings.Join(requirements, ", ")
	}

	return ""
}

func ValidatePasswordMatch(password, confirmPassword string) string {
	if password != confirmPassword {
		return "Passwords do not match"
	}
	return ""
}
