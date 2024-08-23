package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

)

func main() {
	fmt.Println("Running terraform code")
	err := os.Chdir("AKS/")
	if err != nil {
		log.Fatalf("Error changing directory: %v", err)
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

	
	applyCmd := exec.Command("terraform", "apply", "-auto-approve")

	applyCmd.Stdout = os.Stdout
	applyCmd.Stderr = os.Stderr

	err = applyCmd.Run()
	if err != nil {
		log.Fatal("error running Terraform apply:", err)
	}

}
