package server

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"os/exec"
	"strings"
	"time"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// buildv1beta1 "github.com/shipwright-io/build/pkg/apis/build/v1beta1"
)

func CreateBuild() error {
    // Check if build exists and delete if it does
    if err := deleteBuildIfExists(cfg.BuildName); err != nil {
        return fmt.Errorf("error handling existing build: %v", err)
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
    contextDir: %s
  strategy:
    name: %s
    kind: ClusterBuildStrategy
  output:
    image: docker.io/%s/%s:%s
    pushSecret: push-secret
  timeout: %s
`, cfg.BuildName, cfg.GithubURl,cfg.BuildDir ,cfg.BuildStrategy, cfg.RegistryUser, cfg.ImageName, cfg.ImgTag, cfg.Timeout)

    fmt.Println(buildYAML)

    if err := applyYAML(buildYAML); err != nil {
        return fmt.Errorf("error applying build.yaml: %v", err)
    }

    fmt.Println("BUILD APPLIED")

    fmt.Println("build created!!!")
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

func deleteBuildIfExists(buildName string) error {
	checkCmd := exec.Command("kubectl", "get", "build", buildName)
	if _, err := checkCmd.CombinedOutput(); err == nil {
		fmt.Println("BUILD EXISTS, DELETING EXISTING BUILD")
		deleteCmd := exec.Command("kubectl", "delete", "build", buildName)
		if output, err := deleteCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("error deleting existing build: %v, output: %s", err, output)
		}
		fmt.Println("Existing build deleted.")
	} else {
		fmt.Println("Build does not exist, proceeding to create a new one.")
	}
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
		return fmt.Errorf("error applying buildrun.yaml: %v", err)
	}

	fmt.Println("buildrun created!!!")
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


