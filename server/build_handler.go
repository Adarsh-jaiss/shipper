package server

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/adarsh-jaiss/shipper/configs"
	"github.com/gofiber/fiber/v3"

	"gopkg.in/yaml.v2"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	buildv1beta1 "github.com/shipwright-io/build/pkg/apis/build/v1beta1"
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

func CreateBuild() error {
	checkCmd := exec.Command("kubectl", "get", "build", cfg.BuildName)
	if _, err := checkCmd.CombinedOutput(); err == nil {
		fmt.Println("BUILD EXISTS, DELETING EXISTING BUILD")
		deleteCmd := exec.Command("kubectl", "delete", "build", cfg.BuildName)
		if output, err := deleteCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("error deleting existing build: %v, output: %s", err, output)
		}
		fmt.Println("Existing build deleted.")
	} else {
		fmt.Println("Build does not exist, proceeding to create a new one.")
	}

	fmt.Println("Creating build..")

	buildYAML := fmt.Sprintf(`
    apiVersion: shipwright.io/v1beta1
    kind: Build
    metadata:
      name: %s
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
    `, cfg.BuildName, cfg.GithubURl, cfg.BuildStrategy, cfg.RegistryUser, cfg.ImageName, cfg.Timeout)

	fmt.Println(buildYAML)

	if err := applyYAML(buildYAML); err != nil {
		fmt.Errorf("error applying build.yaml: %v", err)
		log.Fatal(err)
	}

	fmt.Println("BUILD APPLIED")

	fmt.Println("build created!!!")
	return nil
}

func BuildRun() error {
	fmt.Println("creating buildrun...")
	checkCmd := exec.Command("kubectl", "get", "buildrun", cfg.BuildName)
	if _, err := checkCmd.CombinedOutput(); err == nil {
		fmt.Println("BUILD EXISTS, DELETING EXISTING BUILD")
		deleteCmd := exec.Command("kubectl", "delete", "buildrun", cfg.BuildName)
		if output, err := deleteCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("error deleting existing build: %v, output: %s", err, output)
		}
		fmt.Println("Existing build deleted.")
	} else {
		fmt.Println("Build does not exist, proceeding to create a new one.")
	}

	fmt.Println("creating buildrun..")

	buildRunYAML := fmt.Sprintf(`
    apiVersion: shipwright.io/v1beta1
    kind: BuildRun
    metadata:
      generateName: %s-buildrun-
    spec:
      build:
        name: %s
    `, cfg.BuildName, cfg.BuildName)

	fmt.Println(buildRunYAML)

	if err := applyBuildRunYAML(buildRunYAML); err != nil {
		fmt.Errorf("error applying buildrun.yaml: %v", err)
		log.Fatal(err)
	}

	fmt.Println("buildrun created!!!")

	// // Extract the buildrun name from the YAML or command output (if generated)
	// var brs *buildv1beta1.BuildRunSpec
	// var br *buildv1beta1.BuildRun
	// buildRunName := GetBuildRunName(brs)

	// // Monitor the BuildRun process and log any errors
	// if err := monitorBuildRun(br, buildRunName); err != nil {
	// 	return fmt.Errorf("error during buildrun: %v", err)
	// }

	// Handle build run deletion if specified
	// if cfg.BuildRunDeletion {
	// 	if err := exec.Command("kubectl", "delete", "buildrun", fmt.Sprintf("%s-run", cfg.BuildName)).Run(); err != nil {
	// 		return fmt.Errorf("error deleting buildrun: %v", err)
	// 	}
	// }
	return nil
}

func applyBuildRunYAML(yamlContent string) error {
	// Parse the YAML content
	var obj map[string]interface{}
	if err := yaml.Unmarshal([]byte(yamlContent), &obj); err != nil {
		return fmt.Errorf("error parsing YAML: %v", err)
	}

	// Check for generateName and replace it with a unique name
	metadata, ok := obj["metadata"].(map[string]interface{})
	if ok {
		if generateName, exists := metadata["generateName"]; exists {
			delete(metadata, "generateName")
			uniqueSuffix := strings.ToLower(time.Now().Format("20060102-150405"))
			metadata["name"] = fmt.Sprintf("%s%s", generateName, uniqueSuffix)
		}
	}

	// Marshal the modified YAML back to a string
	modifiedYAML, err := yaml.Marshal(obj)
	if err != nil {
		return fmt.Errorf("error marshaling modified YAML: %v", err)
	}

	// Apply the modified YAML using kubectl create
	cmd := exec.Command("kubectl", "create", "-f", "-")
	cmd.Stdin = bytes.NewReader(modifiedYAML)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error creating BuildRun: %v, output: %s", err, string(output))
	}

	return nil
}

// use server instead of org : when using both docker and quay
func createDockerRegistrySecret(server, user, password, email string) error {
	checkCmd := exec.Command("kubectl", "get", "secret", "push-secret")
	fmt.Println("PAHUCH GYA!!!")
	if _, err := checkCmd.CombinedOutput(); err == nil {
		// Secret exists, delete it
		fmt.Println("SECRET EXIST, DELETING OLD SECRET")
		deleteCmd := exec.Command("kubectl", "delete", "secret", "push-secret")
		if output, err := deleteCmd.CombinedOutput(); err != nil {
			fmt.Printf("error deleting existing Docker registry secret: %v, output: %s\n", err, output)
		}
	} else {
		fmt.Println("push-secret does not exist, proceeding to create a new one.")
	}

	cmd := exec.Command("kubectl", "create", "secret", "docker-registry", "push-secret",
		"--docker-server="+server,
		"--docker-username="+user,
		"--docker-password="+password,
		"--docker-email="+email)
	fmt.Println("Code REACHED HERE")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("CODE FATA")
		return fmt.Errorf("error creating Docker registry secret: %v, output: %s", err, output)
	}
	return nil
}

func applyYAML(yamlContent string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cat <<EOF | kubectl apply -f - --validate=false\n%s\nEOF", yamlContent))
	cmd.Stdin = bytes.NewBufferString(yamlContent)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error applying YAML: %v, output: %s", err, output)
	}
	return nil
}

func GetBuildRunName(br *buildv1beta1.BuildRunSpec) string {
	name := br.BuildName()
	// Only BuildRuns with a ReferencedBuild can actually return a proper Build name
	return name
}

// Monitor and log the BuildRun process
func monitorBuildRun(br *buildv1beta1.BuildRun, buildRunName string) error {
	var buildRunType buildv1beta1.Type
	for {
		// Check the status of the BuildRun
		if br.Status.CompletionTime != nil {
			if br.IsSuccessful() {
				fmt.Printf("BuildRun %s succeeded\n", buildRunName)
				break
			}
			if br.Status.IsFailed(buildRunType) { // Ensure the correct method signature
				return fmt.Errorf("BuildRun %s failed: %v, Reason: %v", buildRunName, br.Status, br.Status.FailureDetails)
			}
		}

		// Sleep before polling again
		time.Sleep(5 * time.Second)
	}

	return nil
}
