package main

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
	"unicode"
)

func getKubernetesClient(inCluster bool, clusterContext string) *kubernetes.Clientset {
	if inCluster == false {
		// Create config reading from kubernetes config file from .kube directory
		kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig}
		configOverrides := &clientcmd.ConfigOverrides{CurrentContext: clusterContext}
		kubconf, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
		if err != nil {
			log.Fatalln(err)
		}
		// creates the clientset
		client, err := kubernetes.NewForConfig(kubconf)
		if err != nil {
			log.Fatalln(err)
		}
		return client
	}
	// creates the in-cluster config
	kubconf, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	client, err := kubernetes.NewForConfig(kubconf)
	// creates the clientset
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

func CreateSecretData(kubernetesClient *kubernetes.Clientset, namespace, secretName string, secretData map[string]string) {
	secret, err := kubernetesClient.CoreV1().Secrets(namespace).Create(context.TODO(), &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		StringData: secretData,
	}, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Secret %s already exists\n", secretName)
	} else {
		fmt.Printf("Secret %s was created on namespace %s\n", secret.Name, secret.Namespace)
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
