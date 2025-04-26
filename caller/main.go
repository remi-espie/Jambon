package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func main() {
	event := loadConfig("EVENT")
	ollamaHost := loadConfig("OLLAMA_HOST")
	piperHost := loadConfig("PIPER_HOST")
	whisperHost := loadConfig("WHISPER_HOST")

	log.Println("Using event:", event)
	ollama := callPrompt(ollamaHost, event)
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
