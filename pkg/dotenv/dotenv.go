package dotenv

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// GetPath returns the value of ENV_PATH or the default value.
func GetPath(defaultPath string) string {
	path := os.Getenv("ENV_PATH")
	if path == "" {
		path = defaultPath
	}
	return path
}

// Load reads values from the given filepath and copies them to dst.
// It returns an error if any variable is missing.
func Load(filepath string, dst *map[string]string) error {
	// Read env file
	envMap, err := godotenv.Read(filepath)
	if err != nil {
		return err
	}

	// Set env values and catch the missing ones
	missingEnv := []string{}
	for k := range *dst {
		if v, ok := envMap[k]; !ok {
			missingEnv = append(missingEnv, k)
		} else {
			(*dst)[k] = v
		}
	}

	// Return an error with a list of missing variables if any
	if len(missingEnv) != 0 {
		missingEnvStr := strings.Join(missingEnv, ",")
		return fmt.Errorf("missing environment variables: %s", missingEnvStr)
	}
	return nil
}
