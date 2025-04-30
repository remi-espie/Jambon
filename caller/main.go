package main

import (
	"flag"
	"fmt"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	"github.com/nmeilick/go-whisper"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func main() {
	event := loadConfig("EVENT")
	ollamaHost := loadConfig("OLLAMA_HOST")
	whisperHost := loadConfig("WHISPER_HOST")

	log.Println("Using event:", event)

	ollama := callPrompt(ollamaHost, event)
	// Local TTS
	speech := htgotts.Speech{Folder: "audio", Language: voices.French, Handler: &handlers.Native{}}
	err := speech.Speak(ollama)
	if err != nil {
		log.Fatal("Error transforming text to speech:", err)
	}

	// Whisper
	whisperClient := whisper.NewClient(whisper.WithBaseURL(whisperHost))
	response, err := whisperClient.TranscribeFile("audio/answer.wav")
	if err != nil {
		log.Fatalf("Error transcribing file: %v", err)
	}

	fmt.Printf("Transcription: %s\n", response.Text)

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
