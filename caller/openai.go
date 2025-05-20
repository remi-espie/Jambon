package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"io"
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

func speak(client openai.Client, text string) (string, error) {
	response, err := client.Audio.Speech.New(context.Background(), openai.AudioSpeechNewParams{
		Input:          text,
		Model:          "speaches-ai/Kokoro-82M-v1.0-ONNX",
		ResponseFormat: openai.AudioSpeechNewParamsResponseFormatMP3,
		Voice:          "af_heart",
	})
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Error closing response body:", err)
		}
	}(response.Body)
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	file, err := os.Create("audio/" + uuid.New().String() + ".mp3")
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Error closing audio file:", err)
		}
	}(file)
	_, err = file.Write(b)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}
