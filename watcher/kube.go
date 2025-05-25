package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getEvent(client *kubernetes.Clientset) watch.Interface {
	events, err := client.CoreV1().Events("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("Error getting events:", err)
	}
	return events
}

func getInClusterConfig() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal("Error getting in-cluster config:", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Error creating Kubernetes clientset:", err)
	}

	return clientset
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

func launchJob(client *kubernetes.Clientset, event corev1.Event, ollamaHost string, whisperHost string) *batchv1.Job {
	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprint("jambon-caller-", uuid.New().String()),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "jambon-caller",
							Image:   "ghcr.io/remi-espie/jambon-caller:feat-ci",
							Command: []string{"./main"},
							Args:    []string{"-event_name", event.Name, "-event_namespace", event.Namespace},
							Env: []corev1.EnvVar{
								{
									Name: "GIT_SSH_KEY",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "git-ssh-secret",
											},
											Key: "key",
										},
									},
								},
								{
									Name: "GITHUB_TOKEN",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "github-api-secret",
											},
											Key: "key",
										},
									},
								},
								{
									Name:  "OLLAMA_HOST",
									Value: ollamaHost,
								},
								{
									Name:  "WHISPER_HOST",
									Value: whisperHost,
								},
							},
							ImagePullPolicy: "Always",
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
