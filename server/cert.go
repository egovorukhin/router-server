package server

type Certificate struct {
	Name       string `yaml:"name"`
	Cert       string `yaml:"cert"`
	Key        string `yaml:"key"`
	ClientCert string `yaml:"clientCert"`
}

type Certificates []Certificate

func (c Certificates) Get(name string) *Certificate {
	for _, cert := range c {
		if cert.Name == name {
			return &cert
		}
	}
	return nil
}
