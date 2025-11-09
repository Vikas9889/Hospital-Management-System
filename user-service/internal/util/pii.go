package util

import "user-service/internal/repository"

// MaskUserForLogs returns a copy of the user with PII masked for safe logging.
func MaskUserForLogs(u *repository.User) repository.User {
	masked := *u
	if len(masked.Phone) > 4 {
		keep := 2
		masked.Phone = "***" + masked.Phone[len(masked.Phone)-keep:]
	} else if len(masked.Phone) > 0 {
		masked.Phone = "***"
	}
	// mask email: keep first char and domain
	at := -1
	for i := 0; i < len(masked.Email); i++ {
		if masked.Email[i] == '@' {
			at = i
			break
		}
	}
	if at > 1 {
		local := masked.Email[:at]
		domain := masked.Email[at:]
		if len(local) > 2 {
			masked.Email = string(local[0]) + "***" + string(local[len(local)-1]) + domain
		} else {
			masked.Email = "***" + domain
		}
	}
	return masked
}
