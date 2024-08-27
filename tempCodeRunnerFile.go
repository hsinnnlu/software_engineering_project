func init() {
	verifyCodes = make(map[string]VerificationCode)
	loginAttempts = make(map[string]int)
	lockoutTime = make(map[string]time.Time)
}