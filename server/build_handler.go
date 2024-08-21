package server

import (
	"fmt"
	"os"
	"log"
	"os/exec"

	"github.com/adarsh-jaiss/shipper/configs"
	"github.com/gofiber/fiber/v3"
)

func BuildHandler(c fiber.Ctx) (error) {
	var cfg configs.Build
	if err := c.Bind().Body(&cfg); err!= nil {
		return fmt.Errorf("error parsing request:%v",err)	
	}

	fmt.Printf("REGISTRY_SERVER: %v\n", cfg.RegistryServer)
	fmt.Printf("REGISTRY_USER : %v\n", cfg.RegistryUser)
	fmt.Printf("REGISTRY_PASSWORD: %v\n", cfg.RegistryPassword)
	fmt.Printf("REGISTRY_EMAIL: %v\n", cfg.RegistryEmail)

	fmt.Printf("Build Name: %v\n",cfg.BuildName)
	fmt.Printf("Source type: %v\n",cfg.SourceType)
	fmt.Printf("Build Run Deletion: %v\n",cfg.BuildRunDeletion)	//makk as true only
	fmt.Printf("Github url: %v\n",cfg.GithubURl)
	fmt.Printf("BuildStartegy: %v\n",cfg.BuildStrategy)


	fmt.Printf("Build Name: %v\n",cfg.BuildName)
	fmt.Printf("Image Name: %v\n",cfg.ImageName)
	fmt.Printf("Timeout: %v\n",cfg.Timeout)	

	
	// Store the values in environment variables or pass them to the Terraform script
	os.Setenv("REGISTRY_SERVER", cfg.RegistryServer)
	os.Setenv("REGISTRY_USER", cfg.RegistryUser)
	os.Setenv("REGISTRY_PASSWORD", cfg.RegistryPassword)
	os.Setenv("REGISTRY_EMAIL", cfg.RegistryEmail)
	
	fmt.Println("Enviroment variables are setuped correctly", "\n")

	script()
	// return c.JSON(cfg)
	return c.JSON(map[string]string{"msg":"container built successful"})
}

func script(){

	// Change to the directory containing the Terraform scripts
	fmt.Println("Running terraform code")
	err := os.Chdir("scripts/AKS")
	if err != nil {
		log.Fatalf("Error changing directory: %v", err)
	}
	

	// Run the Terraform plan command
	initCmd := exec.Command("terraform", "init")
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr
	err = initCmd.Run()
	if err != nil {
		log.Fatalf("Error running Terraform plan: %v", err)
	}

	// Run the Terraform plan command
	planCmd := exec.Command("terraform", "plan")
	planCmd.Stdout = os.Stdout
	planCmd.Stderr = os.Stderr
	err = planCmd.Run()
	if err != nil {
		log.Fatalf("Error running Terraform plan: %v", err)
	}

	// Run the Terraform apply command
	cmd := exec.Command("terraform", "apply", "-auto-approve")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error running Terraform: %v", err)
	}

}