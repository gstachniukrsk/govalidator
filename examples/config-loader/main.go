package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gstachniukrsk/govalidator"
)

// AppConfig represents the application configuration
type AppConfig struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	Logging  LoggingConfig  `json:"logging"`
	Features FeatureFlags   `json:"features"`
}

// ServerConfig represents HTTP server configuration
type ServerConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"readTimeout"`
	WriteTimeout int    `json:"writeTimeout"`
	TLS          struct {
		Enabled  bool   `json:"enabled"`
		CertFile string `json:"certFile,omitempty"`
		KeyFile  string `json:"keyFile,omitempty"`
	} `json:"tls"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Driver          string `json:"driver"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Database        string `json:"database"`
	MaxConnections  int    `json:"maxConnections"`
	MaxIdleConns    int    `json:"maxIdleConns"`
	ConnMaxLifetime int    `json:"connMaxLifetime"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password,omitempty"`
	DB       int    `json:"db"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Output string `json:"output"`
}

// FeatureFlags represents feature flag configuration
type FeatureFlags struct {
	EnableCache       bool `json:"enableCache"`
	EnableRateLimit   bool `json:"enableRateLimit"`
	EnableMetrics     bool `json:"enableMetrics"`
	EnableDebugMode   bool `json:"enableDebugMode"`
	MaintenanceMode   bool `json:"maintenanceMode"`
}

// Define validation schema for the application configuration
var (
	// Host pattern (hostname or IP)
	hostPattern = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)

	// TLS configuration schema
	tlsSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("enabled").
			Required().
			WithValidators(govalidator.IsBooleanValidator),
		govalidator.NewField("certFile").
			Optional().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(1),
			),
		govalidator.NewField("keyFile").
			Optional().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(1),
			),
	).WithExtra(govalidator.ExtraForbid)

	// Server configuration schema
	serverSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("host").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.RegexpValidator(*hostPattern),
			),
		govalidator.NewField("port").
			Required().
			WithValidators(
				govalidator.IsIntegerValidator,
				govalidator.MinFloatValidator(1),
				govalidator.MaxFloatValidator(65535),
			),
		govalidator.NewField("readTimeout").
			Required().
			WithValidators(
				govalidator.IsIntegerValidator,
				govalidator.MinFloatValidator(1),
				govalidator.MaxFloatValidator(300),
			),
		govalidator.NewField("writeTimeout").
			Required().
			WithValidators(
				govalidator.IsIntegerValidator,
				govalidator.MinFloatValidator(1),
				govalidator.MaxFloatValidator(300),
			),
		govalidator.NewField("tls").
			Required().
			WithSchema(tlsSchema),
	).WithExtra(govalidator.ExtraForbid)

	// Database configuration schema
	databaseSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("driver").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.OneOfValidator("postgres", "mysql", "sqlite"),
			),
		govalidator.NewField("host").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(1),
			),
		govalidator.NewField("port").
			Required().
			WithValidators(
				govalidator.IsIntegerValidator,
				govalidator.MinFloatValidator(1),
				govalidator.MaxFloatValidator(65535),
			),
		govalidator.NewField("username").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(1),
			),
		govalidator.NewField("password").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(1),
			),
		govalidator.NewField("database").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(1),
			),
		govalidator.NewField("maxConnections").
			Required().
			WithValidators(
				govalidator.IsIntegerValidator,
				govalidator.MinFloatValidator(1),
				govalidator.MaxFloatValidator(1000),
			),
		govalidator.NewField("maxIdleConns").
			Required().
			WithValidators(
				govalidator.IsIntegerValidator,
				govalidator.MinFloatValidator(1),
				govalidator.MaxFloatValidator(100),
			),
		govalidator.NewField("connMaxLifetime").
			Required().
			WithValidators(
				govalidator.IsIntegerValidator,
				govalidator.MinFloatValidator(0),
			),
	).WithExtra(govalidator.ExtraForbid)

	// Redis configuration schema
	redisSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("host").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(1),
			),
		govalidator.NewField("port").
			Required().
			WithValidators(
				govalidator.IsIntegerValidator,
				govalidator.MinFloatValidator(1),
				govalidator.MaxFloatValidator(65535),
			),
		govalidator.NewField("password").
			Optional().
			WithValidators(govalidator.IsStringValidator),
		govalidator.NewField("db").
			Required().
			WithValidators(
				govalidator.IsIntegerValidator,
				govalidator.MinFloatValidator(0),
				govalidator.MaxFloatValidator(15),
			),
	).WithExtra(govalidator.ExtraForbid)

	// Logging configuration schema
	loggingSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("level").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.OneOfValidator("debug", "info", "warn", "error"),
			),
		govalidator.NewField("format").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.OneOfValidator("json", "text"),
			),
		govalidator.NewField("output").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.OneOfValidator("stdout", "stderr", "file"),
			),
	).WithExtra(govalidator.ExtraForbid)

	// Feature flags schema
	featuresSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("enableCache").
			Required().
			WithValidators(govalidator.IsBooleanValidator),
		govalidator.NewField("enableRateLimit").
			Required().
			WithValidators(govalidator.IsBooleanValidator),
		govalidator.NewField("enableMetrics").
			Required().
			WithValidators(govalidator.IsBooleanValidator),
		govalidator.NewField("enableDebugMode").
			Required().
			WithValidators(govalidator.IsBooleanValidator),
		govalidator.NewField("maintenanceMode").
			Required().
			WithValidators(govalidator.IsBooleanValidator),
	).WithExtra(govalidator.ExtraForbid)

	// Root configuration schema
	configSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("server").
			Required().
			WithSchema(serverSchema),
		govalidator.NewField("database").
			Required().
			WithSchema(databaseSchema),
		govalidator.NewField("redis").
			Required().
			WithSchema(redisSchema),
		govalidator.NewField("logging").
			Required().
			WithSchema(loggingSchema),
		govalidator.NewField("features").
			Required().
			WithSchema(featuresSchema),
	).WithExtra(govalidator.ExtraForbid)
)

// LoadConfig loads and validates configuration from a JSON file
func LoadConfig(filename string) (*AppConfig, error) {
	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON to any
	var rawConfig any
	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Validate the configuration
	valid, errs := configSchema.ValidateWithPresenter(
		context.Background(),
		rawConfig,
		govalidator.PathPresenter("."),
		govalidator.DetailedErrorPresenter(),
	)

	if !valid {
		fmt.Println("\nConfiguration validation errors:")
		for path, messages := range errs {
			fmt.Printf("  %s:\n", path)
			for _, msg := range messages {
				fmt.Printf("    - %s\n", msg)
			}
		}
		return nil, fmt.Errorf("configuration validation failed")
	}

	// Parse into typed struct
	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// createSampleConfig creates sample configuration files
func createSampleConfig(filename string, valid bool) error {
	var config any

	if valid {
		config = map[string]any{
			"server": map[string]any{
				"host":         "localhost",
				"port":         8080,
				"readTimeout":  30,
				"writeTimeout": 30,
				"tls": map[string]any{
					"enabled":  false,
					"certFile": "",
					"keyFile":  "",
				},
			},
			"database": map[string]any{
				"driver":          "postgres",
				"host":            "localhost",
				"port":            5432,
				"username":        "appuser",
				"password":        "secretpassword",
				"database":        "appdb",
				"maxConnections":  25,
				"maxIdleConns":    5,
				"connMaxLifetime": 300,
			},
			"redis": map[string]any{
				"host":     "localhost",
				"port":     6379,
				"password": "",
				"db":       0,
			},
			"logging": map[string]any{
				"level":  "info",
				"format": "json",
				"output": "stdout",
			},
			"features": map[string]any{
				"enableCache":     true,
				"enableRateLimit": true,
				"enableMetrics":   true,
				"enableDebugMode": false,
				"maintenanceMode": false,
			},
		}
	} else {
		// Invalid configuration for testing
		config = map[string]any{
			"server": map[string]any{
				"host":         "localhost",
				"port":         99999, // Invalid: port out of range
				"readTimeout":  30,
				"writeTimeout": 30,
				"tls": map[string]any{
					"enabled": "yes", // Invalid: should be boolean
				},
			},
			"database": map[string]any{
				"driver":   "mongodb", // Invalid: not in allowed list
				"host":     "localhost",
				"port":     5432,
				"username": "", // Invalid: empty username
				"password": "secretpassword",
				"database": "appdb",
			},
			"redis": map[string]any{
				"host": "localhost",
				"port": 6379,
				"db":   20, // Invalid: Redis DB must be 0-15
			},
			"logging": map[string]any{
				"level":  "trace", // Invalid: not in allowed list
				"format": "json",
				"output": "stdout",
			},
			// Missing "features" section (required)
		}
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func main() {
	validConfigFile := "config.json"
	invalidConfigFile := "config-invalid.json"

	// Create sample config files
	fmt.Println("Creating sample configuration files...")
	if err := createSampleConfig(validConfigFile, true); err != nil {
		log.Fatalf("Failed to create valid config: %v", err)
	}
	fmt.Printf("✓ Created valid config: %s\n", validConfigFile)

	if err := createSampleConfig(invalidConfigFile, false); err != nil {
		log.Fatalf("Failed to create invalid config: %v", err)
	}
	fmt.Printf("✓ Created invalid config: %s\n", invalidConfigFile)

	// Test valid configuration
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TESTING VALID CONFIGURATION")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	config, err := LoadConfig(validConfigFile)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
	} else {
		fmt.Println("✓ Configuration loaded and validated successfully!")
		fmt.Printf("\nServer will run on: %s:%d\n", config.Server.Host, config.Server.Port)
		fmt.Printf("Database: %s on %s:%d\n", config.Database.Driver, config.Database.Host, config.Database.Port)
		fmt.Printf("Redis: %s:%d (DB %d)\n", config.Redis.Host, config.Redis.Port, config.Redis.DB)
		fmt.Printf("Logging: %s level, %s format\n", config.Logging.Level, config.Logging.Format)
		fmt.Printf("Feature Flags: Cache=%v, RateLimit=%v, Metrics=%v\n",
			config.Features.EnableCache,
			config.Features.EnableRateLimit,
			config.Features.EnableMetrics,
		)
	}

	// Test invalid configuration
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TESTING INVALID CONFIGURATION")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	_, err = LoadConfig(invalidConfigFile)
	if err != nil {
		fmt.Println("\n✓ Invalid configuration correctly rejected")
	} else {
		fmt.Println("\n✗ Invalid configuration was incorrectly accepted!")
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("Validation complete!")
	fmt.Printf("Valid config saved to: %s\n", validConfigFile)
	fmt.Printf("Invalid config saved to: %s\n", invalidConfigFile)
	fmt.Println(strings.Repeat("=", 80))
}
