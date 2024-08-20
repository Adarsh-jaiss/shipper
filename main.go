package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	// "sync"

	"fmt"

	"github.com/adarsh-jaiss/shipper/api"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Main() {
	runAPIServer()
	runTerraform()
}

func main() {
	Main()
	// wg := sync.WaitGroup{}
	// wg.Add(2)

	// go func() {
	// 	defer wg.Done()
	// 	runAPIServer()
	// }()

	// go func() {
	// 	defer wg.Done()
	// 	runTerraform()
	// }()

	
	// wg.Wait()

}

func runAPIServer() {
	fmt.Println("Running API server")
	http.HandleFunc("/register-docker-secret", api.RegisterDockerSecret)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error running API server: %v", err)
	}
}

func runTerraform() {
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

	fmt.Println("Terraform script executed successfully")
}
