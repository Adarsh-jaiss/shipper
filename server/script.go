package server

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/adarsh-jaiss/shipper/configs"
)

func Script() error {
	fmt.Println("Running terraform code")
	err := os.Chdir("scripts/AKS")
	if err != nil {
		log.Fatalf("Error changing directory: %v", err)
	}

	clusterExists, err := checkExistingCluster("shipper-resource-group", "shipper-aks-cluster")
	if err != nil {
		return fmt.Errorf("error checking for existing cluster: %v", err)
	}

	err = writeTerraformVarsFile(*cfg, clusterExists)
	if err != nil {
		log.Fatalf("Error writing Terraform variables: %v", err)
	}

	initCmd := exec.Command("terraform", "init")
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr
	err = initCmd.Run()
	if err != nil {
		log.Fatalf("Error running Terraform init: %v", err)
	}

	planCmd := exec.Command("terraform", "plan")
	planCmd.Stdout = os.Stdout
	planCmd.Stderr = os.Stderr
	err = planCmd.Run()
	if err != nil {
		log.Fatalf("Error running Terraform plan: %v", err)
	}

	var applyCmd *exec.Cmd
	if clusterExists {
		applyCmd = exec.Command("terraform", "apply", "-var", "reuse_existing_resources=true", "-auto-approve")
	} else {
		applyCmd = exec.Command("terraform", "apply", "-auto-approve")
	}
	applyCmd.Stdout = os.Stdout
	applyCmd.Stderr = os.Stderr

	err = applyCmd.Run()
	if err != nil {
		return fmt.Errorf("error running Terraform apply: %v", err)
	}

	return nil
}

func checkExistingCluster(resourceGroup, clusterName string) (bool, error) {
	cmd := exec.Command("az", "aks", "show", "--resource-group", resourceGroup, "--name", clusterName)
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Exit status 3 means the resource was not found
			if exitError.ExitCode() == 3 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func writeTerraformVarsFile(cfg configs.Build, reuseExistingResources bool) error {
	existingContent, err := os.ReadFile("terraform.tfvars")
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading terraform.tfvars file: %v", err)
	}

	contentStr := string(existingContent)

	contentStr = updateOrAddTerraformVariable(contentStr, "REGISTRY_SERVER", cfg.RegistryServer)
	contentStr = updateOrAddTerraformVariable(contentStr, "REGISTRY_USER", cfg.RegistryUser)
	contentStr = updateOrAddTerraformVariable(contentStr, "REGISTRY_PASSWORD", cfg.RegistryPassword)
	contentStr = updateOrAddTerraformVariable(contentStr, "REGISTRY_EMAIL", cfg.RegistryEmail)
	contentStr = updateOrAddTerraformVariable(contentStr, "reuse_existing_resources", fmt.Sprintf("%v", reuseExistingResources))

	return os.WriteFile("terraform.tfvars", []byte(contentStr), 0644)
}

func updateOrAddTerraformVariable(content, key, value string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?m)^\s*%s\s*=\s*".*"$`, key))
	content = re.ReplaceAllString(content, "")
	content += fmt.Sprintf(`%s = "%s"`+"\n", key, value)
	return content
}

// func ensureShipwrightInstalled() error {
// 	cmd := exec.Command("kubectl", "get", "deployment", "-n", "shipwright-build", "shipwright-build-controller")
// 	if err := cmd.Run(); err != nil {
// 		return installShipwright()
// 	}
// 	return nil
// }

// func installShipwright() error {
// 	cmd := exec.Command("kubectl", "apply", "--filename", "https://github.com/shipwright-io/build/releases/download/v0.13.0/release.yaml")
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	return cmd.Run()

// }
