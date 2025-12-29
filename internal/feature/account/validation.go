package account

import "regexp"

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func isValidEmail(email string) bool {
	return emailPattern.MatchString(email)
}

func isValidRegistrationPassword(password string) bool {
	return len(password) >= 6
}
