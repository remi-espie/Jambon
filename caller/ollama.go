package main

import (
	"github.com/xyproto/ollamaclient"
	"log"
)

func callPrompt(ollamaHost string, event string) string {
	oc := ollamaclient.New()
	oc.Verbose = true
	oc.API = ollamaHost
	oc.Model = "gemma3"
	if err := oc.PullIfNeeded(); err != nil {
		log.Fatal("Could not pull model:", err)
	}
	prompt := "You are a Kubernetes expert bot, you will have a call with a user tasked to solve the following issue:`" + event + "`. Please provide a detailed response to the user, including the steps they should take to resolve the issue. Be sure to include any relevant commands or configurations that may be helpful. The user is looking for a solution that is easy to understand and implement. Please keep your response concise and focused on the task at hand."
	output, err := oc.GetOutput(prompt)
	if err != nil {
		log.Fatal("Error getting ollama output:", err)
	}
	return output
}
