package configs

type Build struct {
	RegistryServer   string `json:"registryServer"`
	RegistryUser     string `json:"registryUser"`
	RegistryPassword string `json:"registryPassword"`
	RegistryEmail    string `json:"registryEmail"`
	BuildName        string `json:"buildName"`
	// SourceType       string `json:"sourceType"`
	ImgTag           string `json:"imgTag"`
	BuildDir         string `json:"buildDir"`
	GithubURl        string `json:"githubUrl"`
	BuildStrategy    string `json:"buildStrategy"`
	ImageName        string `json:"imageName"`
	Timeout          string `json:"timeout"`
}