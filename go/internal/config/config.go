package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type AWSCredentials struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type AWSConfig struct {
	Region      string         `json:"region"`
	Endpoint    string         `json:"endpoint"`
	Credentials AWSCredentials `json:"credentials"`
}

type S3Config struct {
	BucketName     string `json:"bucket_name"`
	ForcePathStyle bool   `json:"force_path_style"`
}

type SQSConfig struct {
	QueueName string `json:"queue_name"`
}

type DatabaseConfig struct {
	PostgresDSN string `json:"postgres_dsn"`
}

type JWTConfig struct {
	Secret string `json:"secret"`
}

type ServerInfo struct {
	Port string `json:"port"`
	Name string `json:"name"`
}

type ServersConfig struct {
	GoAPI     ServerInfo `json:"go_api"`
	PythonAPI ServerInfo `json:"python_api"`
	ReactApp  ServerInfo `json:"react_app"`
}

type Config struct {
	AWS         AWSConfig      `json:"aws"`
	S3          S3Config       `json:"s3"`
	SQS         SQSConfig      `json:"sqs"`
	Database    DatabaseConfig `json:"database"`
	JWT         JWTConfig      `json:"jwt"`
	Servers     ServersConfig  `json:"servers"`
	Environment string         `json:"environment"`
}

func LoadConfig() *Config {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}

	configPath := filepath.Join(workDir, "..", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "config.json"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("Failed to read config file:", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal("Failed to parse config JSON:", err)
	}

	config = processEnvVars(config)

	if config.Database.PostgresDSN == "" {
		log.Fatal("POSTGRES_DSN environment variable is required")
	}

	return &config
}

func processEnvVars(config Config) Config {
	config.Database.PostgresDSN = resolveEnvVar(config.Database.PostgresDSN)
	
	config.JWT.Secret = resolveEnvVar(config.JWT.Secret)
	
	return config
}

func resolveEnvVar(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		varPart := value[2 : len(value)-1] 
		
		if strings.Contains(varPart, ":") {
			parts := strings.SplitN(varPart, ":", 2)
			envVar := parts[0]
			defaultVal := parts[1]
			
			if envValue := os.Getenv(envVar); envValue != "" {
				return envValue
			}
			return defaultVal
		}
		
		envVar := varPart
		if envValue := os.Getenv(envVar); envValue != "" {
			return envValue
		}
		return "" 
	}
	
	return value
}