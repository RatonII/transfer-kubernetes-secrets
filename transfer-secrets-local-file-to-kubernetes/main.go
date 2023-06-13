package main

import (
	"flag"
	"fmt"
	vault "github.com/sosedoff/ansible-vault-go"
	"log"
	"os"
	"strconv"
)

func main() {
	kubernetesContext := flag.String("k8s-context", "", "Kubernetes context set when outside of cluster")
	namespace := flag.String("namespace", os.Getenv("KUBE_NAMESPACE"), "The namespace where to update the certificate secrets")
	secretFolderPath := flag.String("secret-folder", "", "The path of the folder which contains the services with the secrets")
	ansibleDecryptKey := flag.String("ansible-decrypt-key", "", "The ansible decryption key used to decrypt the encrypted file")
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
	entries, err := os.ReadDir(*secretFolderPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		// The folder path it is using the windows path syntax to change it to mac or linux replace '\\' with '/'
		serviceDir := fmt.Sprintf("%s\\%s", *secretFolderPath, e.Name())
		if dir, err := isDirectory(serviceDir); dir {
			if err != nil {
				panic(err)
			}
			// The folder path it is using the windows path syntax to change it to mac or linux replace '\\' with '/'
			secretFiles, err := os.ReadDir(fmt.Sprintf("%s\\stage", serviceDir))
			if err != nil {
				log.Println(err)
			}
			secrets := map[string]string{}
			for _, secretFile := range secretFiles {
				if isUpper(secretFile.Name()) {
					// The folder path it is using the windows path syntax to change it to mac or linux replace '\\' with '/'
					str, err := vault.DecryptFile(fmt.Sprintf("%s\\stage\\%s", serviceDir, secretFile.Name()), *ansibleDecryptKey)
					if err != nil {
						fmt.Println(err)
						continue
					}
					secrets[secretFile.Name()] = str

				}
			}
			CreateSecretData(kubernetesClient, *namespace, e.Name(), secrets)
		}
	}
}
