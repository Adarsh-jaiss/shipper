package api

import (
	"fmt"
	"net/http"
	"os"
)

func RegisterDockerSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	registryServer := r.FormValue("registry_server")
	registryUser := r.FormValue("registry_user")
	registryPassword := r.FormValue("registry_password")
	registryEmail := r.FormValue("registry_email")

	fmt.Printf("REGISTRY_SERVER: %v", registryServer)
	fmt.Printf("REGISTRY_USER : %v", registryUser)
	fmt.Printf("REGISTRY_PASSWORD: %v", registryPassword)
	fmt.Printf("REGISTRY_EMAIL: %v", registryEmail)

	// Store the values in environment variables or pass them to the Terraform script
	os.Setenv("REGISTRY_SERVER", registryServer)
	os.Setenv("REGISTRY_USER", registryUser)
	os.Setenv("REGISTRY_PASSWORD", registryPassword)
	os.Setenv("REGISTRY_EMAIL", registryEmail)

	fmt.Fprintf(w, "Docker registry secret registered successfully")
}
