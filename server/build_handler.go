package server

import (
	
	"fmt"
	"os"

	"github.com/adarsh-jaiss/shipper/configs"
	"github.com/gofiber/fiber/v3"
)

var cfg *configs.Build

func BuildHandler(c fiber.Ctx) error {
	if err := c.Bind().Body(&cfg); err != nil {
		return fmt.Errorf("error parsing request:%v", err)
	}

	fmt.Printf("REGISTRY_SERVER: %v\n", cfg.RegistryServer)
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

	if err := authenticate(); err != nil {
		return fmt.Errorf("error authenticating: %v", err)
	}
	fmt.Println("Connected to AKS cluster")

	fmt.Println("creating push secret..")
	if err := createDockerRegistrySecret(cfg.RegistryServer, cfg.RegistryUser, cfg.RegistryPassword, cfg.RegistryEmail); err != nil {
		return fmt.Errorf("error creating Docker registry secret: %v", err)
		
	}
	fmt.Println("push secret created!!")

	if err := CreateBuild(); err != nil {
		return fmt.Errorf("error creating build: %v", err)
	}

	if err := BuildRun(); err != nil {
		return fmt.Errorf("error creating buildrun: %v", err)
	
	}

	os.Remove(KubeconfigPath)

	return c.JSON(map[string]string{"msg": "container built successfully"})
}

