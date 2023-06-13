package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

func main() {
	kubernetesContext := flag.String("k8s-context", "", "Kubernetes context set when outside of cluster")
	namespace := flag.String("namespace", os.Getenv("KUBE_NAMESPACE"), "The namespace where to update the certificate secrets")
	secretName := flag.String("secret-name", "", "The source secretname to get data")
	awsSecretARN := flag.String("aws-secret-arn", "", "The destination aws secret ARN where to transfer k8s secret")
	flag.Parse()
	authInCluster := false
	if os.Getenv("AUTH_IN_CLUSTER") != "" {
		auth, err := strconv.ParseBool(os.Getenv("AUTH_IN_CLUSTER"))
		if err != nil {
			log.Fatalln(err)
		}
		authInCluster = auth
	}
	kubernetesClient := getKubernetesClient(authInCluster, *kubernetesContext)
	updateSecretManagerSecret(getSecretManagerClient(), *awsSecretARN, GetSecretData(kubernetesClient, *namespace, *secretName))
}
