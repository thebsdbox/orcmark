package orchestrator

// This will take a standardised service spec and deploy to a kubernetes cluster

import (
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//InvokeKubernetes - this requires a working kubectl
func (s *Service) InvokeKubernetes() error {
	// use the current context in kubeconfig
	var kubeConfigPath string
	//TODO - Allow a configurable path to a kubeconfig
	if home := homeDir(); home != "" {
		kubeConfigPath = filepath.Join(home, ".kube", "config")
	} else {
		return fmt.Errorf("Unable to locate a home directory, or kubeconfig")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	return nil

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
