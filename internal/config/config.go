package config

type Config struct {
	OAuth struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		RedirectURL  string `yaml:"redirect_url"`
		AuthURL      string `yaml:"auth_url"`
		TokenURL     string `yaml:"token_url"`
		UserInfoURL  string `yaml:"userinfo_url"`
	} `yaml:"oauth"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}
