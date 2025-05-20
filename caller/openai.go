package main

import (
	"context"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"log"
	"os"
)

func openaiClient(baseUrl string) openai.Client {
	client := openai.NewClient(option.WithBaseURL(baseUrl + "/v1"))
	return client
}

func transcribeFile(client openai.Client, filePath string) (string, error) {
	audioFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(audioFile *os.File) {
		err := audioFile.Close()
		if err != nil {
			log.Fatal("Error closing audio file:", err)
		}
	}(audioFile)

	response, err := client.Audio.Transcriptions.New(context.Background(), openai.AudioTranscriptionNewParams{
		File:  audioFile,
		Model: "Systran/faster-whisper-small",
	})
	if err != nil {
		return "", err
	}

	return response.Text, nil
}
