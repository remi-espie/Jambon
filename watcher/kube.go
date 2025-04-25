package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func getEvent(client *kubernetes.Clientset) watch.Interface {
	events, err := client.CoreV1().Events("").Watch(context.TODO(), metav1.ListOptions{})
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

func launchJob(client *kubernetes.Clientset, event corev1.Event, ollamaHost string) *batchv1.Job {
	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprint("jambon-caller_", uuid.New().String()),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "jambon-caller",
							Image: "jambon-caller",
							Args:  []string{"-event", event.Message, "-ollama_host", ollamaHost},
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
		},
	}

	create, err := client.BatchV1().Jobs("jambon").Create(context.TODO(), &job, metav1.CreateOptions{})
	if err != nil {
		log.Println("Error creating job:", err)
	}
	return create
}
