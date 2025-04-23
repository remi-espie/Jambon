package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env.example file, proceeding...")
	}

	kubeconfig := loadConfig("kubeconfig")

	fmt.Printf("Using kubeconfig: %s\n", kubeconfig)

	fmt.Println("Hello, World!")
}

func loadConfig(configString string) string {
	config := flag.String(configString, "", fmt.Sprintf("Path to the %s file", configString))
	flag.Parse()

	err := viper.BindEnv(configString, strings.ToUpper(configString))
	if err != nil {
		log.Fatal(err)
	}
	viper.SetDefault(configString, *config)

	return viper.GetString(configString)
}
