package main

import (
	"github.com/xyproto/ollamaclient"
	"log"
	"strings"
)

func createOllamaClient(ollamaHost string) *ollamaclient.Config {
	oc := ollamaclient.New()
	oc.Verbose = true
	oc.API = ollamaHost
	oc.Model = "qwen3:1.7b"

	if err := oc.PullIfNeeded(); err != nil {
		log.Fatal("Could not pull model:", err)
	}
	return oc
}

func initPrompt(ollamaHost string, event string) (*ollamaclient.Config, string) {
	oc := createOllamaClient(ollamaHost)
	prompt := "You are a Kubernetes expert bot, you will have a call with a user tasked to solve the following issue:`" + event + "`. Please provide a concise response to the user, only including the command(s) to fix the issue. The user is looking for a solution that is easy to understand and implement. Please keep your response concise and focused on the task at hand."
	output, err := oc.GetOutputWithSeedAndTemp(prompt, true, 42, 0.7)
	if err != nil {
		log.Fatal("Error getting ollama output:", err)
	}
	split := strings.SplitAfter(output, "</think>")
	return oc, strings.TrimSpace(split[len(split)-1])
}

func answerUser(oc *ollamaclient.Config, message string) string {
	prompt := "The user has received your call. They should have decided if they can intervene or not. Here is their answer: `" + message + "`. If the user can intervene, please answer `true`. Otherwise, answer `false`. Answer only `true` or `false`. Don't add anything else to your answer. This is extremely important."
	output, err := oc.GetOutputWithSeedAndTemp(prompt, true, 42, 0.7)
	if err != nil {
		log.Fatal("Error getting ollama output:", err)
	}
	split := strings.SplitAfter(output, "</think>")
	return strings.TrimSpace(split[len(split)-1])
}
