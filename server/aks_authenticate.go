package server

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func authenticate() error {
	// Replace these variables with your actual values
	

	// Create a credential object using the client secret
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		return fmt.Errorf("failed to obtain a credential: %v", err)
	}

	// Create an AKS client
	aksClient, err := armcontainerservice.NewManagedClustersClient(subscriptionID, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create AKS client: %v", err)
	}
	// Get the cluster admin credentials
	credentialResults, err := aksClient.ListClusterAdminCredentials(context.Background(), resourceGroupName, clusterName, nil)
	if err != nil {
		return fmt.Errorf("failed to get AKS cluster credentials: %v", err)
	}

	if len(credentialResults.Kubeconfigs) == 0 {
		return fmt.Errorf("no kubeconfig found")
	}

	kubeconfig := credentialResults.Kubeconfigs[0].Value

	// Create a unique temporary kubeconfig file to avoid overwriting the default config
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "shipper-aks-config")
	err = os.WriteFile(kubeconfigPath, kubeconfig, 0644)
	if err != nil {
		return fmt.Errorf("failed to write kubeconfig file: %v", err)
	}

	// Create a Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to build config from flags: %v", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	// Example: List all pods in the default namespace
	pods, err := kubeClient.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list pods: %v", err)
	}

	for _, pod := range pods.Items {
		fmt.Printf("Pod Name: %s\n", pod.Name)
	}

	os.Remove(kubeconfigPath)

	return nil
}
