package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
)

var monitoredEnvVars = []string{
	"SERVER_PORT",
	"LOG_LEVEL",
	"KAFKA_BROKERS",
	"KAFKA_OUTPUT_TOPIC",
	"KAFKA_GROUP_ID",
	"SENTRY_DSN",
}

var sensitiveEnvVars = []string{
	"API_KEY",
	"SECRET",
	"PASSWORD",
	"TOKEN",
}

func printEnvs() {
	fmt.Println("=== ENVIRONMENT VARIABLES ===")
	for _, key := range monitoredEnvVars {
		value := os.Getenv(key)

		if isSensitive(key) && value != "" {
			if len(value) > 4 {

				maskedPart := strings.Repeat("*", len(value)-4)
				lastFour := value[len(value)-4:]
				fmt.Printf("%s: %s%s\n", key, maskedPart, lastFour)
			} else {

				fmt.Printf("%s: %s\n", key, strings.Repeat("*", len(value)))
			}
		} else {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	fmt.Println("============================")
}

func isSensitive(key string) bool {
	key = strings.ToUpper(key)
	for _, sensitiveVar := range sensitiveEnvVars {
		if key == sensitiveVar || strings.Contains(key, sensitiveVar) {
			return true
		}
	}
	return false
}

func LoadConfig() error {
	if os.Getenv("ENVIRONMENT") != "" {
		return nil
	}

	curDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		sentry.CaptureException(err)
		return err
	}

	err = godotenv.Load(curDir + "/.env")
	if err != nil {
		log.Printf("Erro ao carregar arquivo .env: %v", err)
		return err
	}

	printEnvs()

	return nil
}

func getEnv[T any](key string, defaultValue T, parser func(string) (T, error)) T {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := parser(value)
	if err != nil {
		log.Printf("Aviso: Valor inválido para %s, usando padrão: %v", key, defaultValue)
		return defaultValue
	}

	return parsed
}

func GetStringEnv(key string, defaultValue string) string {
	return getEnv(key, defaultValue, func(s string) (string, error) { return s, nil })
}

func GetIntEnv(key string, defaultValue int) int {
	return getEnv(key, defaultValue, strconv.Atoi)
}

func GetBoolEnv(key string, defaultValue bool) bool {
	return getEnv(key, defaultValue, strconv.ParseBool)
}

func GetDurationEnv(key string, defaultValue time.Duration) time.Duration {
	return getEnv(key, defaultValue, time.ParseDuration)
}

func GetFloat64Env(key string, defaultValue float64) float64 {
	return getEnv(key, defaultValue, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	})
}

func StringToInt(s string) (int, error) {
	var v int
	if _, err := fmt.Sscanf(s, "%d", &v); err != nil {
		return 0, err
	}
	return v, nil
}

func StringToBool(s string) (bool, error) {
	if v, err := strconv.ParseBool(s); err == nil {
		return v, nil
	}

	switch strings.ToLower(s) {
	case "yes", "y":
		return true, nil
	case "no", "n":
		return false, nil
	default:
		return false, fmt.Errorf("valor booleano inválido: %s", s)
	}
}
