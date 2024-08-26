package server

import (
	
	"fmt"
	"log"
	"os"
	"os/exec"

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
	fmt.Printf("Source type: %v\n", cfg.SourceType)
	// fmt.Printf("Build Run Deletion: %v\n", cfg.BuildRunDeletion)
	fmt.Printf("Github url: %v\n", cfg.GithubURl)
	fmt.Printf("BuildStartegy: %v\n", cfg.BuildStrategy)

	fmt.Printf("Build Name: %v\n", cfg.BuildName)
	fmt.Printf("Image Name: %v\n", cfg.ImageName)
	fmt.Printf("Timeout: %v\n\n", cfg.Timeout)

	if err := authenticate(); err != nil {
		fmt.Errorf("error authenticating: %v", err)
		log.Fatal(err)
	}
	fmt.Println("Connected to AKS cluster")

	testCMD := exec.Command("kubectl", "cluster-info")
	testCMD.Stdout = os.Stdout
	testCMD.Stdin = os.Stdin
	err := testCMD.Run()
	if err != nil {
		log.Fatal("error running Terraform apply:", err)
	}

	testCMD1 := exec.Command("kubectl", "cluster-info")
	testCMD1.Stdout = os.Stdout
	testCMD1.Stdin = os.Stdin
	err = testCMD1.Run()
	if err != nil {
		log.Fatal("error running Terraform apply:", err)
	}

	fmt.Println("creating push secret..")
	if err := createDockerRegistrySecret(cfg.RegistryServer, cfg.RegistryUser, cfg.RegistryPassword, cfg.RegistryEmail); err != nil {
		fmt.Errorf("error creating Docker registry secret: %v", err)
		log.Fatal(err)
	}
	fmt.Println("push secret created!!")

	if err := CreateBuild(); err != nil {
		fmt.Errorf("error creating build: %v", err)
		log.Fatal(err)
	}

	if err := BuildRun(); err != nil {
		fmt.Errorf("error creating buildrun: %v", err)
		log.Fatal(err)
	}

	os.Remove(KubeconfigPath)

	return c.JSON(map[string]string{"msg": "container built successfully"})
}

