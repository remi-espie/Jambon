package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	eventName := loadConfig("EVENT_NAME")
	eventNamespace := loadConfig("EVENT_NAMESPACE")
	ollamaHost := loadConfig("OLLAMA_HOST")
	whisperHost := loadConfig("WHISPER_HOST")
	disableAutofix := loadConfig("AUTOFIX_GET_OUT_OF_MY_WAY")
	sshKey := loadConfig("GIT_SSH_KEY")
	githubToken := loadConfig("GITHUB_TOKEN")

	log.Println("Using event", eventName, "in namespace", eventNamespace)

	clientSet := createK8sClient()
	event := getEvent(clientSet, eventName, eventNamespace)

	autoresolved := false

	// Autofix
	if len(disableAutofix) == 0 {
		repo, repoPath := cloneRepo(sshKey)
		resource := getResourceContents(repoPath)
		oc := createOllamaClient(ollamaHost)

		fixedResource := promptAutofix(oc, event.Message, resource)
		setResourceContents(repoPath, fixedResource)

		commitMessage := promptAutofixCommitMessage(oc, resource, fixedResource)
		pushAutofix(repo, sshKey, commitMessage)
		mergeAutofix(repo, githubToken)

		autoresolved = true
	}

	if !autoresolved {
		startTime := time.Now()
		// AI Clients
		oc := createOllamaClient(ollamaHost)
		whisperClient := openaiClient(whisperHost)

		ollama := promptCallUser(oc, event.Message)
		log.Println("Ollama response:", ollama)

		// TTS
		filepath, err := speak(whisperClient, ollama)
		if err != nil {
			log.Fatal("Error transforming text to speech: ", err)
		}
		log.Println(filepath)

		userTranscription, err := transcribeFile(whisperClient, "audio/random.mp3")
		if err != nil {
			log.Fatalf("Error transcribing file: %v", err)
		}
		fmt.Printf("Transcription: %s\n", userTranscription)
		if userTranscription == "BEEP" {
			log.Println("No transcription available, exiting...")
			os.Exit(0)
		}

		aiAnswer := promptAnswerUser(oc, userTranscription)
		log.Println("Ollama response:", aiAnswer)
		addCSVLine([]string{startTime.Format(time.RFC3339), "N/A", ollama, aiAnswer})

		if aiAnswer == "true" {
			log.Println("User can intervene, exiting...")
		} else {
			log.Fatal("User cannot intervene, exiting...")
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

func createK8sClient() *kubernetes.Clientset {
	clientConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		log.Fatal("Error creating the Kubernetes client config:", err)
	}

	clientSet, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		log.Fatal("Error creating the Kubernetes client set:", err)
	}

	return clientSet
}

func getEvent(clientSet *kubernetes.Clientset, eventName string, eventNamespace string) *corev1.Event {
	event, err := clientSet.CoreV1().Events(eventNamespace).Get(context.TODO(), eventName, metav1.GetOptions{})
	if err != nil {
		log.Fatal("Error getting the event:", err)
	}

	return event
}
