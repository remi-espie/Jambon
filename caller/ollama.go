package main

import (
	"github.com/xyproto/ollamaclient"
	"log"
)

func createOllamaClient(ollamaHost string) *ollamaclient.Config {
	oc := ollamaclient.New()
	oc.Verbose = true
	oc.API = ollamaHost
	oc.Model = "qwen2:0.5b"

	if err := oc.PullIfNeeded(); err != nil {
		log.Fatal("Could not pull model:", err)
	}
	return oc
}

func initPrompt(ollamaHost string, event string) (*ollamaclient.Config, string) {
	oc := createOllamaClient(ollamaHost)
	prompt := "You are a Kubernetes expert bot, you will have a call with a user tasked to solve the following issue:`" + event + "`. Please provide an extremely concise response to the user, only including the command(s) to fix the issue. The user is looking for a solution that is easy to understand and implement. Please keep your response very concise and focused on the task at hand."
	//output, err := oc.GetOutput(prompt)
	output, err := oc.GetOutputWithSeedAndTemp(prompt, true, 42, 0.7)
	if err != nil {
		log.Fatal("Error getting ollama output:", err)
	}
	return oc, output
}

func answerUser(oc *ollamaclient.Config, message string) string {
	prompt := "The user has received your response, and is answering the following: `" + message + "`. Please continue helping him, provide a very concise answer."
	output, err := oc.GetOutputWithSeedAndTemp(prompt, true, 42, 0.7)
	log.Println("Ollama response:", output)
	if err != nil {
		log.Fatal("Error getting ollama output:", err)
	}
	return output
}
