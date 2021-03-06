package orchestrator

// This will take a standardised service spec and deploy to a kubernetes cluster

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//InvokeKubernetes - this requires a working kubectl
func (s *Service) InvokeKubernetes(autoReap bool) error {
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

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployment := s.setKubeSpec()
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		return err
	}
	log.Infof("Created deployment %q.\n", result.GetObjectMeta().GetName())

	// If we're not auto reaping then return once the deployment has been set
	if autoReap == false {
		return nil
	}

	log.Infoln("Will reap the deployment (after 10 seconds) after all replicas are deployed")

	// Wait a minimum of ten seconds before beginning the auto reap
	time.Sleep(10 * time.Second)

	// Retrieve the orcmark deployment
	orcmarkDeployment, err := deploymentsClient.Get("orcmark-deployment", metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Loop over the deployment and watch the replicas until we're at the expected amount
	for orcmarkDeployment.Status.ReadyReplicas != int32(s.Replicas) {
		log.Debugf("Currently [%d] replicas", orcmarkDeployment.Status.Replicas)
		time.Sleep(1 * time.Second)
		// Update the deployment status
		orcmarkDeployment, err = deploymentsClient.Get("orcmark-deployment", metav1.GetOptions{})
		if err != nil {
			return err
		}
	}
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete("orcmark-deployment", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	log.Infof("Succesfully removed the orcmark deployment from Kubernetes")
	return nil
}

// setKubeSpec will build a configuration for the deployment on kubernetes
func (s *Service) setKubeSpec() *appsv1.Deployment {
	replicas := int32(s.Replicas)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "orcmark-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "benchmark",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "benchmark",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    "benchmark",
							Image:   s.Image,
							Command: s.CMD,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
