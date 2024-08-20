package configs

type PushSecretConfig struct {
	RegistryServer   string `json:"registry_server"`
	RegistryUser     string `json:"registry_user"`
	RegistryPassword string `json:"registry_password"`
	RegistryEmail    string `json:"registry_email"`
}

type BuildConfig struct {
	BuildName        string        `json:"buildName"`
	SourceType       SourceType    `json:"sourceType"`
	BuildRunDeletion bool          `json:"buildRunDeletion"`
	GithubURl        string        `json:"githubUrl"`
	BuildStrategy    BuildStrategy `json:"buildStrategy"`
}

type Build struct {
	RegistryServer   string     `json:"registry_server"`
	RegistryUser     string     `json:"registry_user"`
	RegistryPassword string     `json:"registry_password"`
	RegistryEmail    string     `json:"registry_email"`
	BuildName        string     `json:"buildName"`
	SourceType       SourceType `json:"sourceType"`
	BuildRunDeletion bool       `json:"buildRunDeletion"`
	GithubURl        string     `json:"githubUrl"`
	BuildStrategy    string     `json:"buildStrategy"`
	ImageName        string     `json:"imageName"`
	Timeout          string     `json:"timeout"`
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

type BuildRunConfig struct {
	BuildName string `json:"buildName"`
	ImageName string `json:"imageName"`
	Timeout   string `json:"timeout"`
}
