package tools

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// ViperMustGetString is used to verify whether a specific key is
// set in viper; returns error if it is not set.
func ViperMustGetString(key string) (string, error) {
	if !viper.IsSet(key) {
		return "", fmt.Errorf("viper config param not set: %v", key)
	}
	return viper.GetString(key), nil
}

// ViperMustGetInt is used to verify whether a specific key is
// set in viper; returns error if it is not set.
func ViperMustGetInt(key string) (int, error) {
	if !viper.IsSet(key) {
		return 0, fmt.Errorf("viper config param not set: %v", key)
	}
	return viper.GetInt(key), nil
}

// ViperMustGetInt64 is used to verify whether a specific key is
// set in viper; returns error if it is not set.
func ViperMustGetInt64(key string) (int64, error) {
	if !viper.IsSet(key) {
		return 0, fmt.Errorf("viper config param not set: %v", key)
	}
	return viper.GetInt64(key), nil
}

// ViperMustGetBool is used to verify whether a specific key is
// set in viper; returns error if it is not set.
func ViperMustGetBool(key string) (bool, error) {
	if !viper.IsSet(key) {
		return false, fmt.Errorf("viper config param not set: %v", key)
	}
	return viper.GetBool(key), nil
}

// ViperMustGetStringMapString is used to verify whether a specific key is
// set in viper; returns error if it is not set.
func ViperMustGetStringMapString(key string) (map[string]string, error) {
	if !viper.IsSet(key) {
		return map[string]string{}, fmt.Errorf("viper config param not set: %v", key)
	}
	return viper.GetStringMapString(key), nil
}

// ViperMustGetDuration is used to verify whether a specific key is
// set in viper; returns error if it is not set.
func ViperMustGetDuration(key string) (time.Duration, error) {
	if !viper.IsSet(key) {
		return time.Duration(0), fmt.Errorf("viper config param not set: %v", key)
	}
	return viper.GetDuration(key), nil
}
