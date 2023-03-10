package config

type AgonesConfig struct {
	KeyFile    string        `yaml:"keyFile"`
	CertFile   string        `yaml:"certFile"`
	CaCertFile string        `yaml:"caCertFile"`
	Namespace  string        `yaml:"namespace"`
	Allocator  ServerAddress `yaml:"allocator"`
}
