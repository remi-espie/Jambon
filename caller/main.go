package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"log"
	"strings"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	eventName := loadConfig("EVENT_NAME")
	eventNamespace := loadConfig("EVENT_NAMESPACE")
	ollamaHost := loadConfig("OLLAMA_HOST")
	whisperHost := loadConfig("WHISPER_HOST")

	log.Println("Using event", eventName, "in namespace", eventNamespace)

	clientSet := createK8sClient()
	event := getEvent(clientSet, eventName, eventNamespace)

	oc, ollama := initPrompt(ollamaHost, event.Message)
	log.Println("Ollama response:", ollama)

	// Local TTS
	speech := htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.Native{}}
	filepath, err := speech.CreateSpeechFile(ollama, uuid.New().String())
	log.Println(filepath)
	if err != nil {
		log.Fatal("Error transforming text to speech: ", err)
	}

	// Whisper
	whisperClient := openaiClient(whisperHost)
	response, err := transcribeFile(whisperClient, filepath)
	if err != nil {
		log.Fatalf("Error transcribing file: %v", err)
	}
	fmt.Printf("Transcription: %s\n", response)

	answer := answerUser(oc, response)
	filepath, err = speech.CreateSpeechFile(answer, uuid.New().String())
	log.Println(filepath)
	if err != nil {
		log.Fatal("Error transforming text to speech: ", err)
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
