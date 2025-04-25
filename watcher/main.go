package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
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
	if events == nil {
		log.Fatal("No event stream available")
	}
	defer events.Stop()
	eventChannel := events.ResultChan()

	for {
		for event := range eventChannel {
			kubeEvent, ok := event.Object.(*v1.Event)
			if !ok {
				log.Println("Received an object that is not a Kubernetes event")
				continue
			}
			log.Printf("Event type: %s, Name: %s, Reason: %s, Message: %s\n", kubeEvent.Type, kubeEvent.Name, kubeEvent.Reason, kubeEvent.Message)
		}
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
