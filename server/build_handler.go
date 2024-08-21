package server

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

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
	fmt.Printf("Build Run Deletion: %v\n", cfg.BuildRunDeletion) //makk as true only
	fmt.Printf("Github url: %v\n", cfg.GithubURl)
	fmt.Printf("BuildStartegy: %v\n", cfg.BuildStrategy)

	fmt.Printf("Build Name: %v\n", cfg.BuildName)
	fmt.Printf("Image Name: %v\n", cfg.ImageName)
	fmt.Printf("Timeout: %v\n\n", cfg.Timeout)

	if err := script(); err != nil {
		return fmt.Errorf("error running Terraform apply: %v", err)
	}
	// return c.JSON(cfg)
	return c.JSON(map[string]string{"msg": "container built successful"})
}

func writeTerraformVarsFile(cfg configs.Build) error {
    // Read the existing contents of the terraform.tfvars file
    existingContent, err := os.ReadFile("terraform.tfvars")
    if err != nil && !os.IsNotExist(err) {
        return fmt.Errorf("error reading terraform.tfvars file: %v", err)
    }

    // Convert the file content to a string
    contentStr := string(existingContent)

    // Update or add the variables
    contentStr = updateOrAddTerraformVariable(contentStr, "REGISTRY_SERVER", cfg.RegistryServer)
    contentStr = updateOrAddTerraformVariable(contentStr, "REGISTRY_USER", cfg.RegistryUser)
    contentStr = updateOrAddTerraformVariable(contentStr, "REGISTRY_PASSWORD", cfg.RegistryPassword)
    contentStr = updateOrAddTerraformVariable(contentStr, "REGISTRY_EMAIL", cfg.RegistryEmail)

    // Write the updated content to the terraform.tfvars file
    return os.WriteFile("terraform.tfvars", []byte(contentStr), 0644)
}

func updateOrAddTerraformVariable(content, key, value string) string {
    // Regular expression to match all occurrences of the variable in the file
    re := regexp.MustCompile(fmt.Sprintf(`(?m)^\s*%s\s*=\s*".*"$`, key))
    
    // Remove all occurrences of the key
    content = re.ReplaceAllString(content, "")
    
    // Append the new key-value pair
    content += fmt.Sprintf(`%s = "%s"`+"\n", key, value)
    
    return content
}


func script() error {

	// Change to the directory containing the Terraform scripts
	fmt.Println("Running terraform code")
	err := os.Chdir("scripts/AKS")
	if err != nil {
		log.Fatalf("Error changing directory: %v", err)
	}

	// Write the Terraform variables to the terraform.tfvars file
	err = writeTerraformVarsFile(*cfg)
	if err != nil {
		log.Fatalf("Error writing Terraform variables: %v", err)
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

	return cmd.Run()

}
