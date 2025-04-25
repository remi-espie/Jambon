package main

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func getEvent(client *kubernetes.Clientset) *v1.EventList {
	events, err := client.CoreV1().Events("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("Error getting events:", err)
	}
	return events
}

func initKubeClient(kubeconfig *string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal("Error building kubeconfig:", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Error creating Kubernetes clientset:", err)
	}

	return clientset
}
