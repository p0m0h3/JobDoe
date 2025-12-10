package backend

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"path/filepath"

	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sRunner struct {
	Clientset *kubernetes.Clientset
}

func NewK8sRunner() (*K8sRunner, error) {
	clientset, err := NewK8sClientSet()
	if err != nil {
		return nil, err
	}
	return &K8sRunner{
		Clientset: clientset,
	}, nil
}

func NewK8sClientSet() (*kubernetes.Clientset, error) {
	// Try in-cluster config first.
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fallback to kubeconfig if not in a cluster.
		home := homedir.HomeDir()

		var kubeconfig *string
		if home != "" {
			kubeconfig = flag.String("kubeconfig",
				filepath.Join(home, ".kube", "config"),
				"path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "",
				"path to the kubeconfig file")
		}

		flag.Parse()
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	// Create clientset.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func (k *K8sRunner) Run(code string, env map[string]string) (string, error) {
	jobuuid := uuid.NewString()
	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobuuid,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    fmt.Sprintf("%s-worker", jobuuid),
							Image:   "python:3.12-alpine",
							Command: []string{"python", "-c", code},
						},
					},
				},
			},
		},
	}

	k.Clientset.BatchV1().Jobs("default").Create(context.TODO(), &job, metav1.CreateOptions{})
	return jobuuid, nil
}

func (k *K8sRunner) GetOutput(jobuuid string) (string, error) {
	pods, err := k.Clientset.CoreV1().Pods("default").List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: fmt.Sprintf("job-name=%s", jobuuid),
		},
	)
	if err != nil {
		return "", err
	}

	if len(pods.Items) == 0 {
		return "", fmt.Errorf("no pods found")
	}

	pod := pods.Items[0]
	if pod.Status.Phase != corev1.PodSucceeded {
		return "", fmt.Errorf("pod not succeeded")
	}

	req := k.Clientset.CoreV1().Pods("default").GetLogs(
		pod.Name,
		&corev1.PodLogOptions{},
	)

	logStream, err := req.Stream(context.TODO())
	if err != nil {
		return "", err
	}
	defer logStream.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logStream)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
