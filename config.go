package main

import (
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/joho/godotenv"
)

// Store app configuration
// Map to the environment variables sets by env flag
// To skip set the Parameter by environment variable you must sen env flag to "-"
type AppConfig struct {
	// App
	TCPPort            string `env:"PORT"`
	PolybaseUrl        string `env:"POLYBASE_URL"`
	PolybaseKey        string `env:"POLYBASE_KEY"`
	PolybaseCollection string `env:"POLYBASE_COLLECTION"`
	TimelockHost       string `env:"TIMELOCK_HOST"`
	TimelockHash       string `env:"TIMELOCK_HASH"`
	ENSOwnerAdress     string `env:"ENS_OWNER_ADDRESS"`
	ENSOwnerPrivateKey string `env:"ENS_OWNER_PRIVATE_KEY"`
	RpcUrl             string `env:"RPC_URL"`
	EnsMainDomain      string `env:"ENS_MAIN_DOMAIN"`
	EnsResolverAddress string `env:"ENS_RESOLWER_ADDRESS"`
	PinataKey          string `env:"PINATA_KEY"`
}

var conf AppConfig
var once sync.Once

// Get pointer to the app config
// The first run runs setUp function that get configuration values from environ variables
func GetConfig() *AppConfig {
	once.Do(setUp)
	return &conf
}

// Get configuration values from the environ variables
// If ENVIRONMENT variable is dev then trying to get environ variables from the .env file
func setUp() {
	if env := os.Getenv("ENVIRONMENT"); env == "dev" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("No .env files found. Using real environment")
		}

	}
	v := reflect.ValueOf(&conf).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		varName, _ := f.Tag.Lookup("env")
		if varName == "-" {
			continue
		}
		env, ok := os.LookupEnv(varName)
		if ok {
			v.Field(i).SetString(env)
		}

	}
}
