package server

import (
	"bytes"
	"fmt"
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
	fmt.Printf("Build Run Deletion: %v\n", cfg.BuildRunDeletion)
	fmt.Printf("Github url: %v\n", cfg.GithubURl)
	fmt.Printf("BuildStartegy: %v\n", cfg.BuildStrategy)

	fmt.Printf("Build Name: %v\n", cfg.BuildName)
	fmt.Printf("Image Name: %v\n", cfg.ImageName)
	fmt.Printf("Timeout: %v\n\n", cfg.Timeout)

	if err := authenticate(); err != nil {
		return fmt.Errorf("error authenticating: %v", err)
	}

	fmt.Println("Connected to AKS cluster")

	if err := createDockerRegistrySecret(cfg.RegistryServer, cfg.RegistryUser, cfg.RegistryPassword, cfg.RegistryEmail); err != nil {
		return fmt.Errorf("error creating Docker registry secret: %v", err)
	}

	buildYAML := fmt.Sprintf(`
    apiVersion: shipwright.io/v1beta1
    kind: Build
    metadata:
      name: %s
	  annotations:
    	build.build.dev/build-run-deletion: %s
    spec:
      source:
        type: Git
        git:
          url: %s
        contextDir: source-build
      strategy:
        name: %s
        kind: ClusterBuildStrategy
      output:
        image: docker.io/%s/%s:latest
        pushSecret: push-secret
      timeout: %s
    `, cfg.BuildName, cfg.BuildRunDeletion, cfg.GithubURl, cfg.BuildStrategy, cfg.RegistryOrg, cfg.ImageName, cfg.Timeout)

	if err := applyYAML(cfg.RegistryOrg, buildYAML); err != nil {
		return fmt.Errorf("error applying build.yaml: %v", err)
	}

	// Monitor build process
	if err := exec.Command("kubectl", "get", "builds", "-w").Run(); err != nil {
		return fmt.Errorf("error monitoring build process: %v", err)
	}

	buildRunYAML := fmt.Sprintf(`
    apiVersion: shipwright.io/v1beta1
    kind: BuildRun
    metadata:
      generateName: %s-buildrun-
    spec:
      build:
        name: %s
    `, cfg.BuildName, cfg.BuildName)

	if err := applyYAML(cfg.RegistryOrg, buildRunYAML); err != nil {
		return fmt.Errorf("error applying buildrun.yaml: %v", err)
	}

	// Monitor buildrun process
	if err := exec.Command("kubectl", "get", "buildruns", "-w").Run(); err != nil {
		return fmt.Errorf("error monitoring buildrun process: %v", err)
	}

	// Handle build run deletion if specified
	if cfg.BuildRunDeletion == "true" {
		if err := exec.Command("kubectl", "delete", "buildrun", fmt.Sprintf("%s-run", cfg.BuildName)).Run(); err != nil {
			return fmt.Errorf("error deleting buildrun: %v", err)
		}
	}

	return c.JSON(map[string]string{"msg": "container built successfully"})
}

// use server insted of org : when using both docker and quay
func createDockerRegistrySecret(org, user, password, email string) error {
	cmd := exec.Command("kubectl", "create", "secret", "docker-registry", "push-secret",
		"--docker-server=https://index.docker.io/v1/"+org,
		"--docker-username="+user,
		"--docker-password="+password,
		"--docker-email="+email)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error creating Docker registry secret: %v, output: %s", err, output)
	}
	return nil
}

func applyYAML(REGISTRY_ORG, yamlContent string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("REGISTRY_ORG=%s cat <<EOF | kubectl apply -f -\n%s\nEOF", REGISTRY_ORG, yamlContent))
	cmd.Stdin = bytes.NewBufferString(yamlContent)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error applying YAML: %v, output: %s", err, output)
	}
	return nil
}
