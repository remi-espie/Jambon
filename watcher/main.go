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
		log.Println("Error loading .env file, proceeding...")
	}

	kubeconfig := loadConfig("KUBECONFIG")
	ollamaHost := loadConfig("OLLAMA_HOST")
	whisperHost := loadConfig("WHISPER_HOST")
	log.Println("Using kubeconfig:", kubeconfig)
	log.Println("Using ollama host:", ollamaHost)
	log.Println("Using whisper host:", whisperHost)

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
			if kubeEvent.Type == "Warning" {
				launchJob(client, *kubeEvent, ollamaHost, whisperHost)
			}
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
