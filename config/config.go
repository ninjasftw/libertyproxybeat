// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Libertyproxybeat LibertyproxybeatConfig
}

type LibertyproxybeatConfig struct {
	Ssl            SSLConfig            `yaml:"ssl"`
	Period         string               `yaml:"period"`
	URLs           []string             `yaml:"urls"`
	Authentication AuthenticationConfig `yaml:"authentication"`
	Beans          []BeanConfig         `yaml:"beans"`
    Fields         map[string]string          `yaml:"fields"`
}

type SSLConfig struct {
    Cafile string
}

type AuthenticationConfig struct {
	Username string
	Password string
}



type BeanConfig struct {
	Name       string      `yaml:"name"`
	Attributes []Attribute `yaml:"attributes"`
	Keys       []string    `yaml:"keys"`
}

type Attribute struct {
	Name string   `yaml:"name"`
	Keys []string `yaml:"keys"`
}
