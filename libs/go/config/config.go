package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/juankohler/crypto-bot/libs/go/errors"
)

var (
	ErrInvalidEnv = errors.Define("models.invalid_env")
)

const (
	Dev   Env = "DEV"
	Stage Env = "STAGE"
	Prod  Env = "PROD"
)

type Env string

var allowedEnvs = map[string]Env{
	Dev.String():   Dev,
	Stage.String(): Stage,
	Prod.String():  Prod,
}

func NewEnv(env string) (Env, error) {
	if env, ok := allowedEnvs[strings.ToUpper(env)]; ok {
		return env, nil
	}

	return "", errors.New(
		ErrInvalidEnv,
		"invalid env",
		errors.WithMetadata("env", env),
	)
}

func (e Env) String() string {
	return string(e)
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func RequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("missing required environment variable: " + key)
	}

	return value
}

func GetEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func GetEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}
