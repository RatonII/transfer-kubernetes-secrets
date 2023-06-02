package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
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
func getSecretManagerClient() *secretsmanager.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to create AWS SecretManager clietn, %v", err)
	}
	return secretsmanager.NewFromConfig(cfg)
}

func GetSecretData(kubernetesClient *kubernetes.Clientset, namespace, secretName string) map[string]string {
	secret, err := kubernetesClient.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	secretData := map[string]string{}
	for k, v := range secret.Data {
		secretData[k] = string(v)
	}
	return secretData
}

func updateSecretManagerSecret(secretClient *secretsmanager.Client, secretARN string, secretData map[string]string) {
	secretBinary, err := json.Marshal(secretData)
	if err != nil {
		log.Fatalf("failed to encode json, %v", err)
	}
	secret := string(secretBinary)
	resp, err := secretClient.PutSecretValue(context.TODO(), &secretsmanager.PutSecretValueInput{
		SecretId:     &secretARN,
		SecretString: &secret,
	})
	if err != nil {
		log.Fatalf("failed to update secret, %v", err)
	}
	fmt.Println(*resp.Name)
}
