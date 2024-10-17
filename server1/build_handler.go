package server

import (
	"fmt"
	"os"

	"github.com/adarsh-jaiss/shipper/configs"
	"github.com/gofiber/fiber/v3"
	shipclient "github.com/shipwright-io/build/pkg/client/clientset/versioned"
	"k8s.io/client-go/tools/clientcmd"
)

var cfg *configs.Build

func BuildHandler(c fiber.Ctx) error {
	if err := c.Bind().Body(&cfg); err != nil {
		return fmt.Errorf("error parsing request:%v", err)
	}

	fmt.Printf("REGISTRY_SERVER: %v\n", cfg.RegistryServer)
	cfg.RegistryServer = "https://index.docker.io/v1/"
	fmt.Printf("REGISTRY_USER : %v\n", cfg.RegistryUser)
	fmt.Printf("REGISTRY_PASSWORD: %v\n", cfg.RegistryPassword)
	fmt.Printf("REGISTRY_EMAIL: %v\n", cfg.RegistryEmail)

	fmt.Printf("Build Name: %v\n", cfg.BuildName)
	// fmt.Printf("Source type: %v\n", cfg.SourceType)
	// fmt.Printf("Build Run Deletion: %v\n", cfg.BuildRunDeletion)
	fmt.Printf("Github url: %v\n", cfg.GithubURl)
	fmt.Printf("BuildStartegy: %v\n", cfg.BuildStrategy)

	fmt.Printf("Build Name: %v\n", cfg.BuildName)
	fmt.Printf("Image Name: %v\n", cfg.ImageName)
	fmt.Printf("Timeout: %v\n\n", cfg.Timeout)
	cfg.Namespace = "shipper-users"
	fmt.Printf("Namespace: %v\n", cfg.Namespace)

	// // Authenticate and get Kubernetes client
	kubeClient, err := Authenticate()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error authenticating: %v", err)})
	}
	fmt.Println("Connected to AKS cluster")

	// Get the current config to create Shipwright client
	// Use the KubeconfigPath to create the config
	config, err := clientcmd.BuildConfigFromFlags("", KubeconfigPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error building config from kubeconfig: %v", err)})
	}

	// Create Shipwright client
	shipClient, err := shipclient.NewForConfig(config)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error creating Shipwright client: %v", err)})
	}

	fmt.Println("Creating push secret...")
	if err := CreateDockerRegistrySecret(kubeClient, configs.Build{
		RegistryServer:   cfg.RegistryServer,
		RegistryUser:     cfg.RegistryUser,
		RegistryPassword: cfg.RegistryPassword,
		RegistryEmail:    cfg.RegistryEmail,
		Namespace:        cfg.Namespace,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error creating Docker registry secret: %v", err)})
	}
	fmt.Println("Push secret created!")

	if err := CreateBuild(kubeClient, shipClient, configs.Build{
		BuildName:     cfg.BuildName,
		ImgTag:        cfg.ImgTag,
		BuildDir:      cfg.BuildDir,
		GithubURl:     cfg.GithubURl,
		BuildStrategy: cfg.BuildStrategy,
		ImageName:     cfg.ImageName,
		Timeout:       cfg.Timeout,
		Namespace:     cfg.Namespace,
		RegistryUser:  cfg.RegistryUser,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error creating build: %v", err)})
	}

	if err := BuildRun(shipClient, configs.Build{
		BuildName: cfg.BuildName,
		Namespace: cfg.Namespace,
	}); err!= nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error creating buildrun: %v", err)})
	}

	os.Remove(KubeconfigPath)

	return c.JSON(map[string]string{"msg": "container built successfully"})
}
