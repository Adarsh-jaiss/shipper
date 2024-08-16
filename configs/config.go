package configs

type PushSecretConfig struct {
	RegistryUser     string `json:"username"`
	RegistryPassword string `json:"password"`
	OwnerEmail       string `json:"email"`
}

type BuildConfig struct {
	BuildName        string        `json:"buildName"`
	SourceType       SourceType    `json:"sourceType"`
	BuildRunDeletion bool          `json:"buildRunDeletion"`
	GithubURl        string        `json:"githubUrl"`
	BuildStrategy    BuildStrategy `json:"buildStrategy"`
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
	BuildName BuildConfig `json:"buildName"`
	ImageName string `json:"imageName"`
	Timeout string `json:"timeout"`
}

