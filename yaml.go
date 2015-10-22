package swagger2

import (
	"gopkg.in/yaml.v2"
)

// LoadYaml parses the incoming byte array as Swagger 2 YAML data
func LoadYaml(in []byte) (*Swagger, error) {
	var s Swagger
	err := yaml.Unmarshal(in, &s)
	if err != nil {
		return nil, err
	}
	return &s, err
}

// Yaml serializes the Swagger 2 document into YAML format
func (s *Swagger) Yaml() ([]byte, error) {
	return yaml.Marshal(s)
}
