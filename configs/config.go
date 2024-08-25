package configs

type Build struct {
	RegistryServer   string `json:"registryServer"`
	RegistryUser     string `json:"registryUser"`
	RegistryPassword string `json:"registryPassword"`
	RegistryEmail    string `json:"registryEmail"`
	RegistryOrg      string `json:"registryOrg"`
	BuildName        string `json:"buildName"`
	SourceType       string `json:"sourceType"`
	// BuildRunDeletion bool `json:"buildRunDeletion"`
	GithubURl        string `json:"githubUrl"`
	BuildStrategy    string `json:"buildStrategy"`
	ImageName        string `json:"imageName"`
	Timeout          string `json:"timeout"`
}

type SourceType int

const (
	Git SourceType = iota + 1
	Dockerfile
)

type BuildStrategy int

const (
	BuildKit BuildStrategy = iota + 1
	kanio
	Buildah
	Buildpacks_V3
	Buildpacks_V3_Heroku
	Ko
	SourceToImage
)
