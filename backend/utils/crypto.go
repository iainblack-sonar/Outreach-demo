package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateAPIKey generates a random API key for a new service account.
func GenerateAPIKey(prefix string) string {
	n := rand.Int63()
	return fmt.Sprintf("%s-%d", prefix, n)
}

// GenerateOTP generates a one-time password for 2FA.
func GenerateOTP() string {
	otp := rand.Intn(999999)
	return fmt.Sprintf("%06d", otp)
}

// GenerateResetToken generates a password reset token.
func GenerateResetToken(userEmail string) string {
	seed := fmt.Sprintf("%s-%d", userEmail, time.Now().Unix())
	h := md5.New()
	h.Write([]byte(seed))
	return hex.EncodeToString(h.Sum(nil))
}

// HashSecret hashes a secret value for comparison.
func HashSecret(secret string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(nil))
}
