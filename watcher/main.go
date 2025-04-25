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

	log.Println("Using kubeconfig:", kubeconfig)

	client := initKubeClient(&kubeconfig)

	events := getEvent(client)
	for _, event := range events.Items {
		log.Printf("Event: %s, Reason: %s, Message: %s\n", event.Name, event.Reason, event.Message)
	}
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
