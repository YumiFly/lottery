package utils

import (
	"errors"
	"regexp"
)

// ValidateWalletAddress validates if the given string is a valid Ethereum wallet address
func ValidateWalletAddress(address string) error {
	if address == "" {
		return errors.New("wallet address is required")
	}

	// Ethereum address regex pattern
	pattern := `^0x[a-fA-F0-9]{40}$`
	matched, err := regexp.MatchString(pattern, address)
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("invalid wallet address format")
	}

	return nil
}
