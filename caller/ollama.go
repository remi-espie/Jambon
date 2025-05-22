package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/xyproto/ollamaclient"
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

func promptCallUser(oc *ollamaclient.Config, event string) string {
	prompt := "You are a Kubernetes expert bot, you will have a call with a user tasked to solve the following issue:`" + event + "`. Please provide a concise response to the user, only including the command(s) to fix the issue. The user is looking for a solution that is easy to understand and implement. Please keep your response concise and focused on the task at hand."
	output, err := oc.GetOutputWithSeedAndTemp(prompt, true, 42, 0.7)
	if err != nil {
		log.Fatal("Error getting ollama output:", err)
	}
	split := strings.SplitAfter(output, "</think>")
	return strings.TrimSpace(split[len(split)-1])
}

func promptAnswerUser(oc *ollamaclient.Config, message string) string {
	prompt := "The user has received your call. They should have decided if they can intervene or not. Here is their answer: `" + message + "`. If the user can intervene, please answer `true`. Otherwise, answer `false`. Answer only `true` or `false`. Don't add anything else to your answer. This is extremely important."
	output, err := oc.GetOutputWithSeedAndTemp(prompt, true, 42, 0.7)
	if err != nil {
		log.Fatal("Error getting ollama output:", err)
	}
	split := strings.SplitAfter(output, "</think>")
	return strings.TrimSpace(split[len(split)-1])
}

func promptAutofix(oc *ollamaclient.Config, eventMessage string, resource string) string {
	thinkRegex := regexp.MustCompile(`<think>[\s\S]*</think>`)
	contentRegex := regexp.MustCompile(`\x60\x60\x60.*\s([\s\S]*)\x60\x60\x60`)

	prompt := fmt.Sprintf("You are a Kubernetes expert. The following event occurred in the cluster: `%s`.\nThe associated manifest is this one:\n```yaml\n%s\n```\nPlease fix the error. Only tell me the complete, fixed file in a code block. Don't say anything else. Only do modifications that solve the issue, and nothing else. Don't overthink this.", eventMessage, resource)
	output, err := oc.GetOutputWithSeedAndTemp(prompt, true, 42, 0.7)
	log.Println("Ollama response autofix:", output)
	if err != nil {
		log.Fatal("Error getting ollama output:", err)
	}

	output = thinkRegex.ReplaceAllLiteralString(output, "")
	output = contentRegex.FindStringSubmatch(output)[1]

	log.Println("New resource file:", output)
	return output
}
