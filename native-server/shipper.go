package native_server

import (
	"context"
	"fmt"
	"time"

	"github.com/adarsh-jaiss/shipper/configs"
	buildv1beta1 "github.com/shipwright-io/build/pkg/apis/build/v1beta1"
	shipclient "github.com/shipwright-io/build/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateDockerRegistrySecret(kubeClient *kubernetes.Clientset, cfg configs.Build) error {
	secretName := "push-secret"

	// Check if the secret already exists
	_, err := kubeClient.CoreV1().Secrets(cfg.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err == nil {
		fmt.Println("SECRET EXISTS, DELETING OLD SECRET")
		deleteErr := kubeClient.CoreV1().Secrets(cfg.Namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
		if deleteErr != nil {
			return fmt.Errorf("error deleting existing Docker registry secret: %v", deleteErr)
		}
	} else if !errors.IsNotFound(err) {
		return fmt.Errorf("error checking for existing secret: %v", err)
	}

	// Create the new secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: cfg.Namespace,
		},
		Type: corev1.SecretTypeDockerConfigJson,
		StringData: map[string]string{
			".dockerconfigjson": fmt.Sprintf(`{"auths":{"%s":{"username":"%s","password":"%s","email":"%s"}}}`, cfg.RegistryServer, cfg.RegistryUser, cfg.RegistryPassword, cfg.RegistryEmail),
		},
	}

	_, err = kubeClient.CoreV1().Secrets(cfg.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating Docker registry secret: %v", err)
	}

	fmt.Println("Docker registry secret created successfully")
	return nil
}



func CreateBuild(kubeClient *kubernetes.Clientset, shipClient *shipclient.Clientset, cfg configs.Build) error {
	// Check if build exists and delete if it does
	if err := deleteBuildIfExists(shipClient, cfg.BuildName, cfg.Namespace); err != nil {
		return fmt.Errorf("error handling existing build: %v", err)
	}

	fmt.Println("Creating build..")

	// Convert timeout string to time.Duration
	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return fmt.Errorf("invalid timeout format: %v", err)
	}

	var kindtype buildv1beta1.Strategy

	build := &buildv1beta1.Build{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfg.BuildName,
			Namespace: cfg.Namespace,
		},
		Spec: buildv1beta1.BuildSpec{
			Source: &buildv1beta1.Source{
				Type: buildv1beta1.GitType,
				Git: &buildv1beta1.Git{
					URL: cfg.GithubURl,
				},
				ContextDir: &cfg.BuildDir,
			},
			Strategy: buildv1beta1.Strategy{
				Name: cfg.BuildStrategy,
				Kind: kindtype.Kind ,
			},
			Output: buildv1beta1.Image{
				Image: fmt.Sprintf("docker.io/%s/%s:%s", cfg.RegistryUser, cfg.ImageName, cfg.ImgTag),
			},
			Timeout: &metav1.Duration{Duration: timeout},
		},
	}

	_, err = shipClient.ShipwrightV1beta1().Builds(cfg.Namespace).Create(context.TODO(), build, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating build: %v", err)
	}

	fmt.Println("BUILD APPLIED")
	fmt.Println("build created!!!")
	return nil
}

func deleteBuildIfExists(shipClient *shipclient.Clientset, buildName, namespace string) error {
	_, err := shipClient.ShipwrightV1beta1().Builds(namespace).Get(context.TODO(), buildName, metav1.GetOptions{})
	if err == nil {
		fmt.Println("BUILD EXISTS, DELETING EXISTING BUILD")
		deleteErr := shipClient.ShipwrightV1beta1().Builds(namespace).Delete(context.TODO(), buildName, metav1.DeleteOptions{})
		if deleteErr != nil {
			return fmt.Errorf("error deleting existing build: %v", deleteErr)
		}
		fmt.Println("Existing build deleted.")
	} else if !errors.IsNotFound(err) {
		return fmt.Errorf("error checking for existing build: %v", err)
	}
	return nil
}



func BuildRun(shipClient *shipclient.Clientset, cfg configs.Build) error {
	fmt.Println("creating buildrun...")

	// Check if buildrun exists and delete if it does
	if err := deleteBuildRunIfExists(shipClient, cfg.BuildName, cfg.Namespace); err != nil {
		return fmt.Errorf("error handling existing buildrun: %v", err)
	}

	fmt.Println("creating buildrun..")

	buildRun := &buildv1beta1.BuildRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-buildrun-", cfg.BuildName),
			Namespace:    cfg.Namespace,
		},
		Spec: buildv1beta1.BuildRunSpec{
			Build: buildv1beta1.ReferencedBuild{
				Name: &cfg.BuildName,
			},
		},
		
	}

	createdBuildRun, err := shipClient.ShipwrightV1beta1().BuildRuns(cfg.Namespace).Create(context.TODO(), buildRun, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating buildrun: %v", err)
	}

	fmt.Printf("BuildRun created with name: %s\n", createdBuildRun.Name)
	return nil
}

func deleteBuildRunIfExists(shipClient *shipclient.Clientset, buildName, namespace string) error {
	listOpts := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("buildrun.shipwright.io/build.name=%s", buildName),
	}
	buildRuns, err := shipClient.ShipwrightV1beta1().BuildRuns(namespace).List(context.TODO(), listOpts)
	if err != nil {
		return fmt.Errorf("error listing buildruns: %v", err)
	}

	for _, br := range buildRuns.Items {
		fmt.Printf("BUILDRUN EXISTS, DELETING EXISTING BUILDRUN: %s\n", br.Name)
		deleteErr := shipClient.ShipwrightV1beta1().BuildRuns(namespace).Delete(context.TODO(), br.Name, metav1.DeleteOptions{})
		if deleteErr != nil {
			return fmt.Errorf("error deleting existing buildrun: %v", deleteErr)
		}
		fmt.Printf("Existing buildrun %s deleted.\n", br.Name)
	}

	return nil
}